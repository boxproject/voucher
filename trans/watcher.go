// Copyright 2017. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trans

import (
	"context"
	"math/big"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	"github.com/boxproject/voucher/token"
	"github.com/boxproject/voucher/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

var (
	depositEvent  = common.SignEvent("Deposit(address,uint256)")
	withdrawEvent = common.SignEvent("Withdraw(address,uint256)")
	allowEvent    = common.SignEvent("AllowFlow(bytes32)")
	disallowEvent = common.SignEvent("DisallowFlow(bytes32)")
	erc20Event    = common.SignEvent("Transfer(address,address,uint256)")
	walletEvent   = common.SignEvent("WalletS(address)")
)

type EventHandler interface {
	Name() common.Hash
	Scan(log *common.EtherLog) error
}

type State interface {
	GetContractAddress() (common.Address, error)
	SetContractAddress(common.Address) error
	Save(key []byte, value []byte) error
	Load(key []byte) ([]byte, error)
}

type lvldbState struct {
	db localdb.Database
}

func (s *lvldbState) GetContractAddress() (addr common.Address, err error) {
	var bytes []byte
	bytes, err = s.db.Get(config.BANKADDRESS)
	//if err == verrors.NoDataErr {
	//	return addr, nil
	//}
	if err != nil {
		return addr, err
	}
	addr.SetBytes(bytes)

	return addr, nil
}

func (s *lvldbState) SetContractAddress(addr common.Address) (err error) {
	return s.db.Put(config.BANKADDRESS, addr.Bytes())
}

func (s *lvldbState) Save(key []byte, value []byte) error {
	return s.db.Put(key, value)
}

func (s *lvldbState) Load(key []byte) ([]byte, error) {
	return s.db.Get(key)
}

type HandlerContext struct {
	Events          map[common.Hash]EventHandler
	State           State
	EtherURL        string
	Retries         int
	DelayedBlocks   *big.Int
	CursorBlocks    *big.Int
	ContractAddress common.Address
	Db              localdb.Database
	BlockNoFilePath string
	NonceFilePath   string
}

func (hc *HandlerContext) addListener(handlers ...EventHandler) {
	for _, handler := range handlers {
		hc.Events[handler.Name()] = handler
	}
}

type EthHandler struct {
	client         *ethclient.Client
	conf           *HandlerContext
	ethQuitCh      chan struct{}
	quitCh         chan struct{}
	events         map[common.Hash]EventHandler
	accountHandler *AccountHandler
	BtcHandler     *BtcHandler
}

func (w *EthHandler) loadAddr() error {
	caddr, err := w.conf.State.GetContractAddress()
	log.Debug("loadAddr addr:", caddr.Hex())
	if err != nil {
		return err
	}
	w.conf.ContractAddress = caddr
	rpcConn, err := rpc.Dial(w.conf.EtherURL)
	if err != nil {
		return err
	}

	w.client = ethclient.NewClient(rpcConn)
	if len(w.conf.Events) == 0 {
		w.conf.addListener(
			newDepositHandler(w.conf, w.client),
			newAllowHandler(w.conf),
			newDisallowHandler(w.conf),
			newWithdrawHandler(w.conf),
			newERC20Handler(w.conf, w.client))
	}
	w.quitCh = make(chan struct{}, 1)
	w.ethQuitCh = make(chan struct{}, 1)
	return nil
}

func (w *EthHandler) Start() error {
	if err := w.loadAddr(); err != nil {
		return err
	}

	go w.asyStart()
	return nil
}

//TODO
func (w *EthHandler) BtcStart() error {
	if err := w.BtcHandler.Start(); err != nil {
		return err
	}
	return nil
}

func (w *EthHandler) BtcStop() error {
	w.BtcHandler.Stop()
	return nil
}

func (w *EthHandler) asyStart() {
	w.reloadBlock() //加载未处理区块
	go w.listen()
	go w.ethChainHandler() //公链操作
	//go w.accountHandler.Start() //生成账号
}

func (w *EthHandler) reloadBlock() error {
	log.Debug("Block file: %s", w.conf.BlockNoFilePath)

	cursorBlkNumber, err := util.ReadNumberFromFile(w.conf.BlockNoFilePath) //block cfg
	if err != nil {
		return err
	}

	// 获取当前节点上最大区块号
	blk, err := w.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	maxBlkNumber := blk.Number()

	nonce, err := w.client.NonceAt(context.Background(), common.HexToAddress(config.RealTimeStatus.Address), maxBlkNumber)
	if err != nil { //设置初始值nonce值
		log.Error("get nonce err: %s", err)
	}
	util.WriteNumberToFile(w.conf.NonceFilePath, big.NewInt(int64(nonce)))

	// 获取向前推N个区块的big.Int值
	checkBefore := w.conf.DelayedBlocks
	log.Debug("before:: max blkNumber: %s, cursor blkNumber: %s", maxBlkNumber.String(), cursorBlkNumber.String())

	for maxBlkNumber.Cmp(cursorBlkNumber) >= int(checkBefore.Int64()) {
		if err = w.checkBlock(new(big.Int).Add(cursorBlkNumber, checkBefore)); err != nil {
			return err
		}
		// 记录下当前的 blocknumber 供恢复用
		util.WriteNumberToFile(w.conf.BlockNoFilePath, cursorBlkNumber)
		cursorBlkNumber = new(big.Int).Add(cursorBlkNumber, big.NewInt(1))
	}

	log.Debug("after:: cursor blkNumber: %s", cursorBlkNumber.String())
	return nil
}

//TODO
// 上公链操作
func (w *EthHandler) ethChainHandler() {
	log.Debug("ethChainHandler start...")
	//处理未完成pending
	loop := true
	for loop {
		select {
		case <-w.ethQuitCh:
			log.Info("PriEthHandler::SendMessage thread exitCh!")
			loop = false
		case data, ok := <-config.Ecr20RecordChan:
			if ok {
				switch data.Type {
				case config.ECR20_TYPE_ALLOW:
					w.ethAllowHandler(data)
				case config.ECR20_TYPE_DISALLOW:
					w.ethDisAllowHandler(data)
				case config.ECR20_TYPE_APPROVE:
					w.ethTransferHandler(data)
				default:
					log.Info("unknow req:%v", data)
				}
			} else {
				log.Error("read from channel failed")
			}
		}
	}
}

//上公链操作-同意
func (w *EthHandler) ethAllowHandler(record *config.Ecr20Record) error {
	log.Debug("ethAllowHandler...hash:%v", record.Hash.Hex())
	util.NoRWMutex.Lock()
	defer util.NoRWMutex.Unlock()
	bank, err := NewBank(w.conf.ContractAddress, w.client)

	if err != nil {
		log.Error("NewBank error:", err)
		return err
	}
	opts := NewKeyedTransactor()
	nonce, err := util.ReadNumberFromFile(w.conf.NonceFilePath) //nonce file
	if err != nil {
		log.Error("read block info err :%s", err)
		return err
	}

	opts.Nonce = nonce
	log.Debug("current nonce :%d", nonce.Int64())
	if tx, err := bank.Allow(opts, record.Hash); err != nil {
		log.Error("bank allow err :%s", err)
	} else {
		log.Info("eth allow :%v", tx.Hash().Hex())
		nonce = nonce.Add(nonce, big.NewInt(config.NONCE_PLUS))
		util.WriteNumberToFile(w.conf.NonceFilePath, nonce)
	}

	return nil
}

//上公链操作-禁用
func (w *EthHandler) ethDisAllowHandler(record *config.Ecr20Record) error {
	log.Debug("ethDisAllowHandler...hash:%v", record.Hash.Hex())
	util.NoRWMutex.Lock()
	defer util.NoRWMutex.Unlock()
	bank, err := NewBank(w.conf.ContractAddress, w.client)

	if err != nil {
		log.Error("NewBank error:", err)
		return err
	}
	opts := NewKeyedTransactor()
	nonce, err := util.ReadNumberFromFile(w.conf.NonceFilePath) //nonce file
	if err != nil {
		log.Error("read block info err :%s", err)
		return err
	}

	opts.Nonce = nonce
	log.Debug("current nonce :%d", nonce.Int64())
	if tx, err := bank.Disallow(opts, record.Hash); err != nil {
		log.Error("bank disallow err :%s", err)
	} else {
		log.Info("eth disallow :%v", tx.Hash().Hex())
		nonce = nonce.Add(nonce, big.NewInt(config.NONCE_PLUS))
		util.WriteNumberToFile(w.conf.NonceFilePath, nonce)
	}
	return nil
}

//上公链操作-转账
func (w *EthHandler) ethTransferHandler(record *config.Ecr20Record) error {
	log.Debug("ethTransferHandler...hash:%v, wdHash:%v", record.Hash.Hex(), record.WdHash.Hex())
	util.NoRWMutex.Lock()
	defer util.NoRWMutex.Unlock()
	bank, err := NewBank(w.conf.ContractAddress, w.client)

	if err != nil {
		log.Error("NewBank error:", err)
		return err
	}
	opts := NewKeyedTransactor()
	nonce, err := util.ReadNumberFromFile(w.conf.NonceFilePath) //nonce file
	if err != nil {
		log.Error("read block info err :%s", err)
		return err
	}

	opts.Nonce = nonce
	log.Debug("current nonce :%d", nonce.Int64())
	if record.Category.Int64() == config.CATEGORY_ETH { //eth
		tx, err := bank.Withdraw(opts, record.To, record.Amount, record.Hash)
		if err != nil {
			log.Error("withDraw error:%s", err)
		} else {
			log.Info("eth transfer tx:%v, wd:%v, to:%v", tx.Hash().Hex(), record.WdHash.Hex(), record.To.Hex())
			nonce = nonce.Add(nonce, big.NewInt(config.NONCE_PLUS))
			util.WriteNumberToFile(w.conf.NonceFilePath, nonce)

			if err := w.conf.Db.Put([]byte(config.WITHDRAW_TX_PRIFIX+tx.Hash().Hex()), []byte(record.WdHash.Hex())); err != nil {
				log.Error("CheckPrivateKey err: %s", err)
			} else {

				config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_WITHDRAW_TX_WEB, WdHash: record.WdHash, TxHash: tx.Hash().Hex()}
			}
		}
	} else { //eth代币
		if token := token.GetTokenByCategory(record.Category.Int64()); token != nil {
			tx, err := bank.TransferERC20(opts, common.HexToAddress(token.ContractAddr), record.To, record.Amount, record.Hash)
			if err != nil {
				log.Error("withDraw error:%s", err)
			} else {
				log.Info("token transfer tx:%v, wd:%v, to:%v", tx.Hash().Hex(), record.WdHash.Hex(), record.To.Hex())
				nonce = nonce.Add(nonce, big.NewInt(config.NONCE_PLUS))
				util.WriteNumberToFile(w.conf.NonceFilePath, nonce)
				if err := w.conf.Db.Put([]byte(config.WITHDRAW_TX_PRIFIX+tx.Hash().Hex()), []byte(record.WdHash.Hex())); err != nil {
					log.Error("CheckPrivateKey err: %s", err)
				} else {
					config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_WITHDRAW_TX_WEB, WdHash: record.WdHash, TxHash: tx.Hash().Hex()}
				}
			}
		} else {
			log.Error("category: %d undifined err", record.Category.Int64())
			return errors.New("category undifined")
		}
	}

	return nil
}

func (w *EthHandler) listen() {
	ch := make(chan *types.Header)
	sid, err := w.client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		log.Error("CANNOT subscribe the block head. cause: %v", err)
		return
	}

	defer sid.Unsubscribe()

	var backoff = common.DefaultBackoff
	// blocked by for loop
	if err = w.recv(sid, ch); err != nil {
		log.Error("receive from public ethereum node failed. cause: %v", err)
		// reconnect after backoff time
		d := backoff.Duration(w.conf.Retries)
		if d > 0 {
			time.Sleep(d)
			go w.listen()
		}
		log.Error("retry %d times and exited, cause: %v", w.conf.Retries, err)
	}
	log.Debug("watcher listener stopped.")
}

func (w *EthHandler) Stop() {
	log.Info("close server...")
	if w.quitCh != nil {
		w.quitCh <- struct{}{}
	}
	if w.ethQuitCh != nil {
		w.ethQuitCh <- struct{}{}
	}
	//if w.accountHandler.quitChannel != nil {
	//	w.accountHandler.Stop() //停止生成账号
	//}

	//w.BtcHandler.Stop()
	log.Info("close server OK")
}

//recv 公链请求
func (w *EthHandler) recv(sid ethereum.Subscription, ch <-chan *types.Header) error {
	var err error
	for {
		select {
		case <-w.quitCh:
			return nil
		case err = <-sid.Err():
			if err != nil {
				return err
			}
		case head := <-ch:
			if head.Number == nil {
				continue
			}
			if err = w.checkBlock(head.Number); err != nil {
				return err
			}
		}
	}
}

func (w *EthHandler) checkBlock(blkNumber *big.Int) error {
	checkPoint := new(big.Int).Sub(blkNumber, w.conf.DelayedBlocks)

	//log.Debug("[HEADER] FromBlock : %s, blkNumber: %s", checkPoint.String(), blkNumber.String())
	log.Debug("[HEADER] blkNumber: %s, checkpoint: %s", blkNumber.String(), checkPoint.String())

	if logs, err := w.client.FilterLogs(
		context.Background(),
		ethereum.FilterQuery{
			FromBlock: checkPoint,
			ToBlock:   checkPoint,
		}); err != nil {
		log.Debug("FilterLogs :%s", err)
		return err
	} else {
		if len(logs) != 0 {
			for _, enventLog := range logs {
				log.Debug("enventLog: %s", enventLog.Topics[0].Hex())
				handler, ok := w.conf.Events[enventLog.Topics[0]]
				log.Debug("enventLog ok: %s", ok)
				if !ok {
					continue
				}
				if err = handler.Scan(&enventLog); err != nil {
					return err
				}
			}
		}
		util.WriteNumberToFile(w.conf.BlockNoFilePath, checkPoint)
	}

	return nil
}

func NewEthHandler(cfg *config.Config, db localdb.Database) (*EthHandler, error) {
	state := &lvldbState{db}
	conf := &HandlerContext{
		Events:          make(map[common.Hash]EventHandler),
		State:           state,
		EtherURL:        cfg.EthereumConfig.Scheme,
		Retries:         cfg.EthereumConfig.Retries,
		DelayedBlocks:   big.NewInt(int64(cfg.EthereumConfig.DelayedBlocks)),
		CursorBlocks:    big.NewInt(int64(cfg.EthereumConfig.CursorBlocks)),
		Db:              db,
		BlockNoFilePath: cfg.EthereumConfig.BlockNoFilePath,
		NonceFilePath:   cfg.EthereumConfig.NonceFilePath,
	}
	w := &EthHandler{conf: conf, quitCh: make(chan struct{}, 1)}
	w.BtcHandler, _ = NewBtcHandler(cfg, db)
	//w.accountHandler = NewAccountHandler(db, w, cfg.EthereumConfig.AccountPoolSize, cfg.EthereumConfig.NonceFilePath)
	return w, nil
}

//发布bank
func (e *EthHandler) DeployBank() (addr common.Address, tx *types.Transaction, b *Bank, err error) {
	log.Debug("DeployBank..........")
	rpcConn, err := rpc.Dial(e.conf.EtherURL)
	if err != nil {
		return
	}
	defer rpcConn.Close()
	client := ethclient.NewClient(rpcConn)
	opts := NewKeyedTransactor()
	if addr, tx, b, err = DeployBank(opts, client); err == nil {
		e.conf.State.SetContractAddress(addr) //地址写入db
	} else {
		log.Error("deploy bank failed. cause: %v", err)
	}
	return
}

//发布地址合约
func (e *EthHandler) DeployWallet(nonce *big.Int) (addr common.Address, err error) {
	log.Debug("DeployWallet..........")
	rpcConn, err := rpc.Dial(e.conf.EtherURL)
	if err != nil {
		return
	}
	defer rpcConn.Close()
	client := ethclient.NewClient(rpcConn)
	opts := NewKeyedTransactor()
	opts.Nonce = nonce
	if bankAddress, err := e.GetContractAddr(); err != nil {
		log.Error("get bankAddress failed. cause: %v", err)
		return addr, err
	} else {
		addr, tx, _, err := DeployWallet(opts, client, bankAddress)
		if err != nil {
			log.Error("deploy wallet failed. cause: %v", err)
			return addr, err
		} else {
			log.Debug("Successful deploy wallet addr: %s", addr.Hex())
			log.Debug("Successful deploy wallet tx: %s", tx.Hash().Hex())
		}
		return addr, nil
	}
}

//获取合约地址
func (e *EthHandler) GetContractAddr() (common.Address, error) {
	return e.conf.State.GetContractAddress()
}
