package trans

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	"github.com/boxproject/voucher/util"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"encoding/hex"
	rpcclient "github.com/boxproject/lib-bitcore/serpcclient"
	"github.com/boxproject/voucher/common"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

//BTC db 前缀
const (
	BTC_LOG_LABLE = "[BTC]"
	BTC_TXID      = "BTC_TXID_"
	BTC_TXID_0    = BTC_TXID + "0_"
	BTC_TXID_1    = BTC_TXID + "1_"
)

const (
	BTC_TXID_TYPE_DEP = "0" //充值
	BTC_TXID_TYPE_WD  = "1" //提现
	BTC_TXID_TYPE_OT  = "2" //其他
)

//BTC
type BTCTxidInfo struct {
	WDhash common.Hash
	TXID   string
	Vout   uint32
	Addr   string
	Type   string //"0"-充值，"1"提现,"2"其他
	Amount *big.Int
}

type BtcHandler struct {
	db        localdb.Database
	btcConf   *config.BitcoinConfig
	clientCfg *rpcclient.ConnConfig
	client    *rpcclient.Client
	accHandle *AccountHandlerBtc
	isStart   bool
	Net       *chaincfg.Params
}

func NewBtcHandler(cfg *config.Config, db localdb.Database) (*BtcHandler, error) {
	btcConf := &cfg.BitcoinConfig
	btcConf.BlockNoFilePath = btcConf.BlockNoFilePath
	clientDfg := &rpcclient.ConnConfig{
		Host:         btcConf.Host,
		User:         btcConf.Rpcuser,
		Pass:         btcConf.Rpcpass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	btcHandler := &BtcHandler{
		db:        db,
		btcConf:   btcConf,
		clientCfg: clientDfg,
		isStart:   false}
	btcHandler.accHandle = NewAccountHandlerBtc(db, btcHandler)

	return btcHandler, nil
}

func (w *BtcHandler) Status() (isWork bool) {
	if w.isStart && w.accHandle.isInited {
		isWork = true
	} else {
		isWork = false
	}
	return isWork
}
func (w *BtcHandler) Start() (err error) {
	if w.isStart {
		log.Info("bitcoin started")
		return nil
	}
	log.Info("bitcoin rpc connect...")
	if w.client, err = rpcclient.New(w.clientCfg, nil); err != nil {
		log.Error("bitcoin rpc new failed.")
		return err
	}
	chainInfo, err := w.client.GetBlockChainInfo()
	if err != nil {
		log.Error("bitcoin rpc conect:", err)
	}

	var netType string
	if w.btcConf.Type != "" {
		netType = w.btcConf.Type
	} else if chainInfo != nil {
		netType = chainInfo.Chain
	} else {
		netType = "main"
	}

	switch netType {
	case "main":
		w.Net = &chaincfg.MainNetParams
		break
	case "test":
		w.Net = &chaincfg.TestNet3Params
		break
	case "regtest":
		w.Net = &chaincfg.RegressionNetParams
		break
	}
	log.Info("bitcoin rpc connect OK")

	if err := w.accHandle.Init(); err != nil {
		return err
	}
	go w.btcChainHandler() //bitcoin公链操作

	return nil
}
//获取btc地址
func (w *BtcHandler) GetAccount() string {
	if w.accHandle.Account != nil {
		return w.accHandle.Account.EncodeAddress()
	}
	return ""
}

func (w *BtcHandler) getTxidAddr(txid string, vout uint32) (addr string, err error) {
	sTxid, _ := chainhash.NewHashFromStr(txid)
	tx, err := w.client.GetTransaction(sTxid)
	if err != nil {
		//log.Println(err)
		return "", err
	}
	bytes, _ := hex.DecodeString(tx.Hex)
	if raw, err := w.client.DecodeRawTransaction(bytes); err != nil {
		//log.Println(err)
		return "", err
	} else {
		return raw.Vout[vout].ScriptPubKey.Addresses[0], nil
	}
}

func (w *BtcHandler) btcDispoitScan(txid string) {
	if sTxid, err := chainhash.NewHashFromStr(txid); err != nil {
		log.Error("get txid hash failed:", err)
	} else {
		var vinAddr = make(map[string]string)
		var fromAddr string
		tx, err := w.client.GetTransaction(sTxid)
		if err != nil {
			log.Error(err)
			return
		}
		bytes, _ := hex.DecodeString(tx.Hex)
		if raw, err := w.client.DecodeRawTransaction(bytes); err != nil {
			log.Error(err)
			return
		} else {
			//获取 from 地址
			fromAddr = ""
			for _, VinValue := range raw.Vin {
				addr, _ := w.getTxidAddr(VinValue.Txid, VinValue.Vout)
				if _, ok := vinAddr[addr]; ok == false {
					vinAddr[addr] = addr
					if fromAddr == "" {
						fromAddr = addr
					} else {
						fromAddr = fromAddr + "," + addr
					}
				}
			}
			for _, value := range raw.Vout {
				if len(value.ScriptPubKey.Addresses) <= 0 {
					continue
				}
				//判断记录内是否有系统内的账户
				if w.accHandle.Account.EncodeAddress() == value.ScriptPubKey.Addresses[0] {
					address := value.ScriptPubKey.Addresses[0]
					amount, _ := btcutil.NewAmount(value.Value)
					var newAmount *big.Int
					newAmount = big.NewInt(int64(amount.ToBTC() * 1e8))

					w.accHandle.uncfmMu.Lock()
					w.accHandle.UncfmTxidMap[txid] = &BTCTxidInfo{TXID: txid, Addr: address, Type: BTC_TXID_TYPE_DEP, Vout: value.N, Amount: newAmount} //Vout暂时设置为0
					w.accHandle.uncfmMu.Unlock()

					if bytes, err := json.Marshal(w.accHandle.UncfmTxidMap[txid]); err != nil {
						log.Error("db unmarshal err: %v", err)
					} else {
						log.Debug("[发现充值txid]:", w.accHandle.UncfmTxidMap[txid])
						//添加记录到数据文件
						if err := w.db.Put([]byte(BTC_TXID_0+txid), bytes); err != nil {
							log.Error("write btc txid txid db failed....")
						}
					}
				}
			}
		}
	}
}

func (w *BtcHandler) btcGetFromAddr(txid *chainhash.Hash) string {
	//获取from地址
	var fromAddr string = ""
	tx, err := w.client.GetTransaction(txid)
	if err != nil {
		log.Error(err)
		return fromAddr
	}
	bytes, err := hex.DecodeString(tx.Hex)
	if err != nil {
		log.Error(err)
		return fromAddr
	}
	raw, err := w.client.DecodeRawTransaction(bytes)
	if err != nil {
		log.Error(err)
		return fromAddr
	}
	var vinAddr = make(map[string]string)
	//获取 from 地址
	for _, VinValue := range raw.Vin {
		addr, _ := w.getTxidAddr(VinValue.Txid, VinValue.Vout)
		if _, ok := vinAddr[addr]; ok == false {
			vinAddr[addr] = addr
			if fromAddr == "" {
				fromAddr = addr
			} else {
				fromAddr = fromAddr + "," + addr
			}
		}
	}
	return fromAddr
}

// 上公链操作
func (w *BtcHandler) btcChainHandler() {
	log.Debug("btcChainHandler start...")
	/*
		功能:
			定时扫描
			1，监控bitcoin节点的充值情况
			2，监控bitcoin节点交易【提币】情况
	*/
	timerListenBtc := time.NewTicker(time.Second * 5)
	isScaning := false

	w.isStart = true
	for w.isStart {
		select {
		case data, ok := <-config.BtcRecordChan: //提币交易
			if ok {
				log.Debug("BtcRecordChannel :%s", data)
				switch data.Type {
				case config.BTC_TYPE_APPROVE:
					if err := w.btcTransferHandler(data); err != nil {
						//如果构造交易失败，返回为空的txid
						config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_WITHDRAW_TX_WEB, To: data.To, WdHash: data.WdHash, TxHash: "", Category: big.NewInt(config.CATEGORY_BTC)}
					}
				default:
					log.Info("unknow req:%v", data)
				}
			} else {
				log.Error("read from channel failed")
			}
		case <-timerListenBtc.C: //BTC公链区块扫描
			if isScaning == true || w.accHandle.isInited == false {
				continue
			}
			isScaning = true
			blockCount, err := w.client.GetBlockCount()
			if err != nil {
				log.Error(err)
				isScaning = false
				continue
			}
			var blockCorsor int64 = 0
			ret, err := util.ReadNumberFromFile(w.btcConf.BlockNoFilePath)
			if err != nil {
				log.Error("get btc block height num failed,", err)
				util.WriteNumberToFile(w.btcConf.BlockNoFilePath, big.NewInt(blockCount))
				blockCorsor = blockCount
			} else {
				blockCorsor = ret.Int64()
			}

			if blockCorsor > blockCount {
				//如果当前的区块高度小于区块的检索高度，设置区块检索高度为当前的区块高度
				util.WriteNumberToFile(w.btcConf.BlockNoFilePath, big.NewInt(blockCount))
				blockCorsor = blockCount
			}

			//如果最新区块高度大于已经扫描了的区块高度
			for blockCount > blockCorsor {
				log.Info("BTC block scanning, height----->", blockCorsor+1)
				if chanHash, err := w.client.GetBlockHash(blockCorsor + 1); err != nil {
					log.Error("getblockhash err:", err)
				} else {
					if msgBlock, err := w.client.GetBlock(chanHash); err != nil {
						log.Error("getblock err:", err)
					} else {
						//检测充值交易
						for _, txid := range msgBlock.Transactions {
							//判断是否为发起的提现交易
							if _, isHaveTxid := w.accHandle.UncfmTxidMap[txid.TxHash().String()]; isHaveTxid == false {
								//检索，查找充值
								w.btcDispoitScan(txid.TxHash().String())
							}
						}
						//检测记录的交易是否已经生效【暂定6次确定】
						w.accHandle.uncfmMu.Lock()
						for txid, txidInfo := range w.accHandle.UncfmTxidMap {
							if sTxid, err := chainhash.NewHashFromStr(txid); err != nil {
								log.Error("get txid hash failed:", err)
							} else {
								if tranResult, err := w.client.GetTransaction(sTxid); err != nil {
									log.Error("get Transaction failed:", err)
								} else {
									//确认区块儿超过6个
									if tranResult.Confirmations >= w.btcConf.Confirmations {
										switch txidInfo.Type {
										case BTC_TXID_TYPE_WD:
											grpcSend := &config.GrpcStream{Type: config.GRPC_WITHDRAW_WEB, To: txidInfo.Addr, Amount: txidInfo.Amount, TxHash: txid, WdHash: txidInfo.WDhash, Category: big.NewInt(config.CATEGORY_BTC)}
											config.ReportedChan <- grpcSend
											//上报提现状态
											log.Info("[WITHDRAW REPORT]:提现上报")
											log.Debug(grpcSend)
											break
										case BTC_TXID_TYPE_DEP:
											//获取from地址
											fromAddr := w.btcGetFromAddr(sTxid)
											grpcSend := &config.GrpcStream{Type: config.GRPC_DEPOSIT_WEB, From: fromAddr, To: txidInfo.Addr, Amount: txidInfo.Amount, TxHash: txid, Category: big.NewInt(config.CATEGORY_BTC)}
											config.ReportedChan <- grpcSend
											//充值确认
											log.Info("[DEPOSIT REPORT]:充值上报", grpcSend)
											log.Debug(grpcSend)
											break
										default:
											log.Info("no this type:", txidInfo.Type)

										}

										if bytes, err := json.Marshal(w.accHandle.UncfmTxidMap[txid]); err != nil {
											log.Error("db unmarshal err: %v", err)
										} else {
											delete(w.accHandle.UncfmTxidMap, txid)
											//跟新数据库
											w.db.Delete([]byte(BTC_TXID_0 + txid))
											if err := w.db.Put([]byte(BTC_TXID_1+txid), bytes); err != nil {
												log.Error("CheckPrivateKey err: %s", err)
											}
										}
									}
								}
							}
						}
						w.accHandle.uncfmMu.Unlock()

						//扫描完毕，更新数据
						util.WriteNumberToFile(w.btcConf.BlockNoFilePath, big.NewInt(blockCorsor+1))
						blockCorsor++
						log.Debug("blockCorsor wirte:", blockCorsor)
					}
				}
			}
			isScaning = false
		default:

		}
	}

	log.Debug("end of btcChainHandler...")
}

// 提款交易
func (w *BtcHandler) btcTransferHandler(record *config.BtcRecord) error {
	//更换地址格式
	pubKeyToAddr, err := btcutil.DecodeAddress(record.To, w.Net)
	if err != nil {
		log.Error("decode address error:", err)
		return err
	} else {
		log.Info("send to addr:", pubKeyToAddr)
	}

	//每千字节的收费(0.0001BTC/kB),计算公式: (148 * 输入数额) + (34 * 输出数额) + 10,没1000字节的费用默认是0.0001BTC
	//handlefee=((148*len(vin) + 34*len(vout) +10)/1000 + 1) * 0.0001BTC
	getAmount := btcutil.Amount(0)                         //能获取的数目
	handlefee := btcutil.Amount(record.Handlefee.Int64())  //手续费
	realOutAmount := btcutil.Amount(record.Amount.Int64()) //要花费的数目
	allSpend := realOutAmount + handlefee
	leftAmount := btcutil.Amount(-1) //剩余找零
	type prvIndex struct {
		pkScript []byte
		//index    uint32
	}
	//[txid]prvIndex, 用于记录txid对应的地址索引
	var txidMap = make(map[wire.OutPoint]*prvIndex)
	var rawtx = wire.NewMsgTx(2)
	////设置txout
	pkScript, err := txscript.PayToAddrScript(pubKeyToAddr)
	if err != nil {
		log.Error(err)
		return err
	}
	txOut := wire.NewTxOut(int64(realOutAmount), pkScript)
	rawtx.AddTxOut(txOut)

	//地址格式转化
	listUspent, err := w.client.ListUnspentMinMaxAddresses(0, 99999999, []btcutil.Address{w.accHandle.Account})
	if err != nil {
		log.Error("get listuspent err:", err)
	}
	for _, _utxo := range listUspent {
		if int64(_utxo.Amount*1e8) <= 0 {
			continue
		}
		haveAmount, _ := btcutil.NewAmount(_utxo.Amount)
		getAmount = getAmount + haveAmount

		//构造交易 txin
		prehash, _ := chainhash.NewHashFromStr(_utxo.TxID)
		prevOut := wire.NewOutPoint(prehash, _utxo.Vout)
		txIn := wire.NewTxIn(prevOut, nil, nil)
		rawtx.AddTxIn(txIn)

		////设置txout
		txinPkScript, err := hex.DecodeString(_utxo.ScriptPubKey)
		if err != nil {
			log.Error(err)
			return err
		}
		txidMap[*prevOut] = &prvIndex{pkScript: txinPkScript}

		if getAmount > allSpend {
			leftAmount = getAmount - allSpend
			//设置找零
			var reedAddr btcutil.Address
			reedAddr, err = btcutil.DecodeAddress(_utxo.Address, w.Net)
			pkScript, err := txscript.PayToAddrScript(reedAddr)
			if err != nil {
				log.Error(err)
				return err
			}
			txOut := wire.NewTxOut(int64(leftAmount), pkScript)
			rawtx.AddTxOut(txOut)
			break
		}
	}

	if leftAmount < 0 {
		log.Error("utxo not enough...")
		return err
	}

	for i, txid := range rawtx.TxIn {
		log.Debug("index:", i, " value:", txid)

		sigScript, err := txscript.SignTxOutput(w.Net, rawtx, i, txidMap[txid.PreviousOutPoint].pkScript,
			txscript.SigHashAll, GetBtcPrivKey(), nil, nil)
		if err != nil {
			fmt.Println(err)
		}
		rawtx.TxIn[i].SignatureScript = sigScript

		//check
		vm, err := txscript.NewEngine(txidMap[txid.PreviousOutPoint].pkScript, rawtx, i,
			txscript.StandardVerifyFlags, nil, nil, -1)
		if err != nil {
			log.Error(err)
			return err
		}
		if err := vm.Execute(); err != nil {
			log.Error(err)
			return err
		} else {
			log.Info("Transaction successfully signed")
		}
	}

	signedTransaction := rawtx
	// Publish the signed sweep transaction.
	txHash, err := w.client.SendRawTransaction(signedTransaction, false)
	if err != nil {
		log.Error("Failed to publish transaction:", err)
		return err
	}
	log.Info("send raw transh hash:", txHash)

	//更新未确认列表
	w.accHandle.uncfmMu.Lock()
	w.accHandle.UncfmTxidMap[txHash.String()] = &BTCTxidInfo{WDhash: record.WdHash, TXID: txHash.String(), Addr: record.To, Type: BTC_TXID_TYPE_WD, Vout: 0, Amount: record.Amount} //Vout暂时设置为0
	w.accHandle.uncfmMu.Unlock()

	if bytes, err := json.Marshal(w.accHandle.UncfmTxidMap[txHash.String()]); err != nil {
		log.Error("db unmarshal err: %v", err)
	} else {
		//提现交易txid生成上报
		if err := w.db.Put([]byte(BTC_TXID_0+txHash.String()), bytes); err != nil {
			log.Error("[DB]put txid err: %s", err)
		}

		//config.ReportedChan <- &config.RepM{RepType: config.REP_WITHDRAW_TX, WdHash: record.WdHash.Hex(), TxHash: txHash.String(), To: record.To}
		config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_WITHDRAW_TX_WEB, To: record.To, WdHash: record.WdHash, TxHash: txHash.String(), Category: big.NewInt(config.CATEGORY_BTC)}
		log.Debug("record UncfmTxidMap [key]:", txHash.String())
		log.Debug("record UncfmTxidMap [value]:", record.WdHash.String())
	}

	return nil
}

//import btc address to bitcoincore
func (w *BtcHandler) btcAddressImport(address btcutil.Address) bool {
	if w.client == nil {
		log.Error("btc rpc client is nil")
		return false
	}
	if addrValid, err := w.client.ValidateAddress(address); err == nil {
		if addrValid.IsWatchOnly == false {
			//如果没有导入
			if err := w.client.ImportAddressRescan(address.EncodeAddress(), address.EncodeAddress(), false); err != nil {
				log.Error("import bitcoin address:", err)
				return false
			} else {
				log.Info("import bitcoin address:", address.EncodeAddress())
				return true
			}
		} else {
			return true
		}
	} else {
		return false
	}
}

func (w *BtcHandler) Stop() {
	//log.Info("stop bitcoin server")
	//stop account
	w.accHandle.Stop()

	w.isStart = false
	if w.client != nil {
		w.client.Shutdown()
	}
	log.Info("bitcoin server Stopped")
}

/*
	---------------------账户相关--------------
*/
type AccountHandlerBtc struct {
	db          localdb.Database
	//quitChannel chan int
	isInited    bool

	handler      *BtcHandler
	uncfmMu      sync.Mutex
	UncfmTxidMap map[string]*BTCTxidInfo //未确认的txid

	Account *btcutil.AddressPubKeyHash
}

//初始化
func NewAccountHandlerBtc(db localdb.Database, handler *BtcHandler) *AccountHandlerBtc {
	return &AccountHandlerBtc{
		isInited:     false,
		db:           db,
		handler:      handler,
		//quitChannel:  make(chan int, 1),
		UncfmTxidMap: make(map[string]*BTCTxidInfo),
	}
}

func (a *AccountHandlerBtc) Init() error {
	log.Debug("AccountHandlerBtc init.......")
	if accStr, err := GetBtcPubKey(a.db); err != nil {
		log.Error(err)
		return err
	} else {
		log.Debug("btc db address:", accStr)
		//转化地址格式
		if addr, err := btcutil.DecodeAddress(accStr, &chaincfg.Params{}); err != nil {
			log.Error(err)
			return err
		} else {
			a.Account, err = btcutil.NewAddressPubKeyHash(addr.ScriptAddress(), a.handler.Net)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	log.Info("account:\n", a.Account)
	//获取btc未被确认的txid列表
	if nconfirmedTxidList, err := a.db.GetPrifix([]byte(BTC_TXID_0)); err != nil {
		log.Error("db error:", err)
	} else {
		a.uncfmMu.Lock()
		for _, value := range nconfirmedTxidList {
			txidInfo := &BTCTxidInfo{}
			if err := json.Unmarshal([]byte(value), txidInfo); err != nil {
				log.Error("db unmarshal err: %v", err)
			} else {
				a.UncfmTxidMap[txidInfo.TXID] = txidInfo
				log.Debug("hash:", txidInfo.WDhash, "txid:", txidInfo.TXID, "type:", txidInfo.Type)
			}
		}
		a.uncfmMu.Unlock()
	}
	go func() {
		timerScanAddImp := time.NewTicker(time.Second * 5)
		isScaning := false
		for a.isInited == false {
			select {
				case <-timerScanAddImp.C:
					if isScaning == true {
						continue
					}
					isScaning = true
					//地址导入
					if a.handler.btcAddressImport(a.Account) {
						a.isInited = true
						log.Info("[BTC]ALL ACCOUNT IS IMPORTED")
					}
					isScaning = false
				default:
			}
		}
	}()

	log.Info("AccountHandlerBtc init finished")
	return nil
}

func (a *AccountHandlerBtc) Stop() {
	a.isInited = false
}

func Btc_test(cfg *config.Config, db localdb.Database) (*BtcHandler, error) {
	/*bitcoin test ----begin
	btc,err  := trans.Btc_test(cfg, db )
	if err != nil {
		logger.Error("run Btc_test failed. cause: %v", err)
		return err
	}
	defer btc.Stop()
	<-quitCh
	return err
	bitcoin test ----end*/

	RecoverPrivateKey(db, []byte("dhjflaksjhdfkasdfkahskd"))
	btcHandler, _ := NewBtcHandler(cfg, db)
	btcHandler.Start()

	return btcHandler, nil
}
