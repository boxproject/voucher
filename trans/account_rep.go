package trans

import (
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/util"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
)

type AccountHandler struct {
	db              localdb.Database
	quitChannel     chan int
	ethHandler      *EthHandler
	accountMu       sync.Mutex
	accountPoolSize int
	//blockNoFilePath    string
	nonceFilePath    string
}

//初始化
func NewAccountHandler(db localdb.Database, handler *EthHandler, poolSize int, path string) *AccountHandler {
	return &AccountHandler{db: db, ethHandler: handler, accountPoolSize: poolSize, nonceFilePath: path, quitChannel: make(chan int, 1)}
}

func (a *AccountHandler) use(address common.Address) {
	//db
	defer a.accountMu.Unlock()
	a.accountMu.Lock() //操作map 加锁
	if config.AccountUsedMap[address.Hex()] != "" {
		log.Error("account: %s already used!", address.Hex())
		return
	}
	a.db.Delete([]byte(config.ACCOUNT_NO_USED + address.Hex()))
	a.db.Put([]byte(config.ACCOUNT_USED+address.Hex()), []byte(address.Hex()))
	config.AccountUsedMap[address.Hex()] = address.Hex() //map中记录使用地址
	if addr, err := a.generateAccount(); err != nil {
		log.Error("generateAccount err:", err)
	} else {
		log.Debug("use add addr :", addr)
		//config.ReportedChan <- &config.RepM{RepType: config.REP_ACCOUNT_ADD, Account: addr.Hex(), Category: config.CATEGORY_ETH}
	}
}

//生成合约账户
func (a *AccountHandler) generateAccount() (addr common.Address, err error) {
	log.Debug("generateAccount.......")

	nonce, err := util.ReadNumberFromFile(a.nonceFilePath) //block cfg
	if err != nil {
		log.Error("read block info err :%s", err)
		return
	}

	log.Debug("current nonce :%d", nonce.Int64())
	if addr, err = a.ethHandler.DeployWallet(nonce); err != nil {
		log.Error("deploy err:", err)
	} else {
		nonce = nonce.Add(nonce, big.NewInt(config.NONCE_PLUS))
		util.WriteNumberToFile(a.nonceFilePath, nonce)
		a.db.Put([]byte(config.ACCOUNT_NO_USED+config.ACCOUNT_TYPE_ETH+addr.Hex()), []byte(addr.Hex()))
	}
	return
}

//生成指定数量的account
func (a *AccountHandler) generateAccounts(count int) (addrs []common.Address, err error) {
	var i int
	for i = 0; i < count; i++ {
		var addr common.Address
		if addr, err = a.generateAccount(); err != nil {
			break
		} else {
			addrs = append(addrs, addr)
		}
	}
	return
}

//启动 读取目前池中多少可用账户 生成指定数量账户
func (a *AccountHandler) Start() {
	log.Debug("AccountHandler start.......")

	a.accountMu.Lock() //操作map 加锁
	accountUsedDbMap, err := a.db.GetPrifix([]byte(config.ACCOUNT_USED))
	if err != nil {
		log.Error("db error:%v", err)
	}

	for _, v := range accountUsedDbMap {
		account := string(v)
		config.AccountUsedMap[account] = account
	}

	log.Debug("accountUsedDbMap...", accountUsedDbMap)

	accountNoUsedDbMap, err := a.db.GetPrifix([]byte(config.ACCOUNT_NO_USED))
	noUsedLen := len(accountNoUsedDbMap)
	log.Debug("accountNoUsedDbMap...", accountNoUsedDbMap)
	log.Debug("noUsedLen...", noUsedLen)
	if noUsedLen < a.accountPoolSize {
		if addrs, err := a.generateAccounts(a.accountPoolSize - noUsedLen); err != nil {
			log.Error("generateAccount err:", err)
		} else {
			log.Debug("addrs.....", addrs)
		}
	}

	a.accountMu.Unlock()

	//chan
	loop := true
	for loop {
		select {
		case <-a.quitChannel:
			log.Info("PriEthHandler::SendMessage thread exitCh!")
			loop = false
		case account, ok := <-config.AccountUsedChan:
			if ok {
				a.use(account)
			} else {
				log.Error("read from channel failed")
			}
		}
	}
}

//启动 读取目前池中多少可用账户 生成指定数量账户
func (a *AccountHandler) Stop() {
	a.quitChannel <- 0
}
