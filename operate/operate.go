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

package operate

import (
	"strings"
	"sync"

	"encoding/json"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/config"
	verrors "github.com/boxproject/voucher/errors"
	"github.com/boxproject/voucher/localdb"
	"github.com/boxproject/voucher/token"
	"github.com/boxproject/voucher/trans"
	"github.com/boxproject/voucher/util"
)

//prefix+host ->pass
var passMap map[string]string = make(map[string]string)
var codeMap map[string]string = make(map[string]string)

//prefix+host ->role
//var roleMap map[string]string = make(map[string]string)

//var roleStartMap map[string]string = make(map[string]string)

//首次code
//var firstCode string = ""

var batchCountMap map[string]int = make(map[string]int)

var authorizedMap map[string]bool = make(map[string]bool)

var apiMu sync.Mutex

type OperateHandler struct {
	quitChannel chan int
	cfg         *config.Config
	db          localdb.Database
	ethHander   *trans.EthHandler
}

func InitOperateHandler(cfg *config.Config, ldb localdb.Database, ethHander *trans.EthHandler) *OperateHandler {
	return &OperateHandler{quitChannel: make(chan int, 1), cfg: cfg, db: ldb, ethHander: ethHander}
}

//签名机操作处理
func (handler *OperateHandler) Start() {
	log.Info("OperateHandler start...")
	loop := true
	for loop {
		select {
		case <-handler.quitChannel:
			log.Info("OperateHandler::SendMessage thread exitCh!")
			loop = false
		case data, ok := <-config.OperateChan:
			if ok {
				switch data.Type {
				case config.VOUCHER_OPERATE_ADDKEY:
					handler.addKey(data)
					break
				case config.VOUCHER_OPERATE_CREATE:
					handler.create(data)
					break
				case config.VOUCHER_OPERATE_DEPLOY:
					handler.deploy(data)
					break
				case config.VOUCHER_OPERATE_START:
					handler.start(data)
					break
				case config.VOUCHER_OPERATE_PAUSE:
					handler.pause(data)
					break
				case config.VOUCHER_OPERATE_HASH_ENABLE:
					handler.hashAllow(data)
					break
				case config.VOUCHER_OPERATE_HASH_DISABLE:
					handler.hashDisAllow(data)
					break
				case config.VOUCHER_OPERATE_HASH_LIST:
					handler.hashList(data)
					break
				case config.VOUCHER_OPERATE_TOKEN_ADD:
					handler.addToken(data)
					break
				case config.VOUCHER_OPERATE_TOKEN_DEL:
					handler.delToken(data)
					break
				case config.VOUCHER_OPERATE_TOKEN_LIST:
					handler.tokenList(data)
					break
				case config.VOUCHER_OPERATE_COIN:
					handler.coin(data)
					break
				default:
					log.Info("unknow asy req: %s", data.Type)
				}
			} else {
				log.Error("OperateHandler read from channel failed")
			}
		}
	}
}

//关闭签名机操作处理
func (handler *OperateHandler) Close() {
	if handler.quitChannel != nil {
		close(handler.quitChannel)
	}
	log.Info("OperateHandler closed")
}

func (handler *OperateHandler) addKey(operate *config.Operate) {
	publicKey := operate.PublicKey
	log.Debug("addKey...")

	if dDate := util.CBCDecrypter([]byte(publicKey), util.GetAesKeyRandomFromDb(handler.db)); len(dDate) == 0 { //解密
		log.Error("key decrypter failed.")
		return
	} else {
		//db
		if err := handler.db.Put([]byte(config.APP_KEY_PRIFIX+operate.AppId), dDate); err != nil {
			log.Error("app keystore err: %s", err)
		} else {
			config.RealTimeStatus.KeyStoreStatus = append(config.RealTimeStatus.KeyStoreStatus, config.KeyStoreStatu{ApplyerId: operate.AppId, ApplyerName: operate.AppName})
			landStatus(handler.db, "")
		}
	}
}

func (handler *OperateHandler) hashAllow(operate *config.Operate) {
	hash := operate.Hash
	log.Debug("hashAllow: %s", hash)
	passWord := getKey(operate.Password, operate.Sign, handler.db, operate.AppId)
	if passWord == "" {
		log.Error("password get failed.")
		return
	}
	if config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_STATED { //服务已启动
		b, passWordStatus := checkKey(handler.db, config.REQ_HASH, handler.cfg.Secret.AppNum, operate.AppId, passWord, true)
		config.RealTimeStatus.Status = passWordStatus
		if !b {
			log.Info("check failed")
			config.RealTimeStatus.Status = passWordStatus
		}
		landStatus(handler.db, config.REQ_HASH)
	} else {
		log.Error("hash allow failed: server not started")
	}
}

func (handler *OperateHandler) hashDisAllow(operate *config.Operate) {
	hash := operate.Hash
	log.Debug("hashDisAllow: %s", hash)
	passWord := getKey(operate.Password, operate.Sign, handler.db, operate.AppId)
	if passWord == "" {
		log.Error("password get failed.")
		return
	}
	if config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_STATED { //服务已启动
		b, passWordStatus := checkKey(handler.db, config.REQ_HASH, handler.cfg.Secret.AppNum, operate.AppId, passWord, false)
		config.RealTimeStatus.Status = passWordStatus
		if !b {
			log.Info("check failed")
			config.RealTimeStatus.Status = config.PASSWORD_STATUS_FAILED
		}
		landStatus(handler.db, config.REQ_HASH)
	} else {
		log.Error("hash allow failed: server not started")
	}
}

func (handler *OperateHandler) create(operate *config.Operate) {
	log.Debug("createHandler...")
	passWord := getKey(operate.Password, operate.Sign, handler.db, operate.AppId)
	//passWord := getKey("xLGBJ4JCmTQCVzrIBEnZGw==", "n0OVW9KHnwmGXJ5rCH/eyFFec4EKdgV+bKZT/jkaJj3q/Z2Gw4wAPE5Si7GPh6PSasQ4E0jBb6dlxciQj7PX+buVXgNevzibUF7hXSK3Fnb182oI6H8hevRJD9FDT4T/5/rRJ375P2Z8NzGXDqEfrcij9O3f5g6nhetF+egTbhLph8MvNxdOZ6NDKBOYeYS2oOCVu52ba2bI509w98V7X9MHLV4HPT3bZwtGlRawnaCJ+UDODmw3LdzF3rRlsi447EBpdzWiKvXqAYHy0lcu0eNpPPkc+5NzFvu9AQWdWRgfpXWx/BmZuUun+lhaH7c/8uVS3Urkearcsfy5oiaa+g==", handler.db, operate.AppId)
	if passWord == "" {
		log.Error("password get failed.")
		return
	}
	if config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_UNCREATED { //服务未创建
		if d, err := trans.GetSecret(handler.db); err != nil {
			log.Error("get secret error:%v", err)
		} else {
			config.RealTimeStatus.D = d
			b, passWordStatus := generateKey(handler.db, config.REQ_CREATE, handler.cfg.Secret.AppNum, operate.AppId, passWord, operate.Code, true)
			config.RealTimeStatus.Status = passWordStatus
			if b {
				config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_CREATED
				config.RealTimeStatus.Address = trans.PrivateKeyToHex(handler.db)
				log.Debug("privateKey generated....")
			} else {
				log.Info("create failed: %s")
			}
			landStatus(handler.db, config.REQ_CREATE)
		}
	} else {
		log.Error("cannot create, status :%s", config.RealTimeStatus.Status)
	}
}

func (handler *OperateHandler) deploy(operate *config.Operate) {
	log.Debug("deployHandler...")
	passWord := getKey(operate.Password, operate.Sign, handler.db, operate.AppId)
	if passWord == "" {
		log.Error("password get failed.")
		return
	}
	b, passWordStatus := checkKey(handler.db, config.REQ_DEPLOY, handler.cfg.Secret.AppNum, operate.AppId, passWord, true)
	config.RealTimeStatus.Status = passWordStatus
	if b {
		if contractAddress, txHash, _, err := handler.ethHander.DeployBank(); err != nil { //发布合约
			log.Error("deploy failed error:", err)
		} else {
			config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_DEPLOYED
			config.RealTimeStatus.ContractAddress = contractAddress.Hex()

			log.Info("deploy success contractAddress:%v, txHash:%v", contractAddress.Hash().Hex(), txHash.Hash().Hex())
		}
	} else {
		log.Info("deploy failed!")
	}
	landStatus(handler.db, config.REQ_DEPLOY)
}

func (handler *OperateHandler) start(operate *config.Operate) {
	log.Debug("startHandler...")
	passWord := getKey(operate.Password, operate.Sign, handler.db, operate.AppId)
	if passWord == "" {
		log.Error("password get failed.")
		return
	}
	if config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_DEPLOYED || config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_PAUSED { //服务已发布或停止
		b, passWordStatus := checkKey(handler.db, config.REQ_START, handler.cfg.Secret.AppNum, operate.AppId, passWord, true)
		config.RealTimeStatus.Status = passWordStatus
		if b {
			if err := handler.ethHander.Start(); err == verrors.NoDataErr {
			} else if err != nil {
				config.RealTimeStatus.Status = config.PASSWORD_SYSTEM_FAILED
				log.Error("start err:%v", err)
			} else {
				config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_STATED
				for _, coinStatu := range config.RealTimeStatus.CoinStatus {
					if coinStatu.Category == config.CATEGORY_BTC && coinStatu.Used == true {
						if err := handler.ethHander.BtcStart(); err != nil {
							log.Error("btc start failed. err: %s", err)
						} else {
							log.Info("btc start success. ")
						}
					}

				}
				log.Info("start success!")
			}
		} else {
			log.Info("start failed.")
		}
		landStatus(handler.db, config.REQ_START)
	} else {
		log.Error("cannot start, status :%s", config.RealTimeStatus.Status)
	}
}

func (handler *OperateHandler) pause(operate *config.Operate) {
	log.Debug("pauseHandler...")
	if config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_STATED { //服务未启动
		//_, b := checkKey(handler.db, config.REQ_PAUSE, handler.cfg.Secret.AppNum, operate.AppId, operate.Password, true)
		//if b {
		handler.ethHander.Stop()
		config.RealTimeStatus.Status = config.PASSWORD_STATUS_OK
		config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_PAUSED
		log.Info("stop success!")
		landStatus(handler.db, config.REQ_PAUSE)
	}
	//} else {
	//	config.RealTimeStatus.Status = config.PASSWORD_STATUS_FAILED
	//
	//	log.Info("pause failed.")
	//}
	//} else {
	//	log.Error("cannot deploy, status :%s", config.RealTimeStatus.Status)
	//}
}

func (handler *OperateHandler) hashList(operate *config.Operate) {
	log.Debug("hashList...")
	config.RealTimeStatus.Status = config.PASSWORD_STATUS_OK
	config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_PAUSED
	landStatus(handler.db, config.REQ_PAUSE)
}

func (handler *OperateHandler) addToken(operate *config.Operate) {
	log.Debug("addToken...")
	if !checkSign(operate.ContractAddr, operate.Sign, operate.AppId, handler.db) {
		log.Info("token add sign check err")
		return
	}
	tokenInfo := &config.TokenInfo{TokenName: operate.TokenName, Decimals: operate.Decimals, ContractAddr: operate.ContractAddr}

	if b, l := token.AddTokenMap(tokenInfo, handler.db); !b {
		log.Info("add token failed.")
		return
	} else {
		config.RealTimeStatus.TokenCount = l
	}
	landStatus(handler.db, "")
}

func (handler *OperateHandler) delToken(operate *config.Operate) {
	log.Debug("delToken...")
	if !checkSign(operate.ContractAddr, operate.Sign, operate.AppId, handler.db) {
		log.Info("token add sign check err")
		return
	}
	tokenInfo := &config.TokenInfo{TokenName: operate.TokenName, Decimals: operate.Decimals, ContractAddr: operate.ContractAddr}

	if b, l := token.DelTokenMap(tokenInfo, handler.db); !b {
		log.Info("add token failed.")
	} else {
		config.RealTimeStatus.TokenCount = l
	}
	landStatus(handler.db, "")
}

func (handler *OperateHandler) tokenList(operate *config.Operate) {
	log.Debug("tokenList...")
	token.TokenList(handler.db)
}

func (handler *OperateHandler) coin(operate *config.Operate) {
	log.Debug("coin...")

	opSuccess := false
	if operate.CoinCategory == config.CATEGORY_BTC {
		if operate.CoinUsed {
			if err := handler.ethHander.BtcStart(); err != nil {
				log.Error("btc start failed. err: %s", err)
			} else {
				opSuccess = true
				config.RealTimeStatus.BtcAddress = handler.ethHander.BtcHandler.GetAccount()
				log.Info("btc start success.")
			}
		} else {
			if err := handler.ethHander.BtcStop(); err != nil {
				log.Error("btc stop failed. err: %s", err)
			} else {
				opSuccess = true
				log.Info("btc stop success.")
			}
		}
	}
	if opSuccess {
		var newCoinStatus []config.CoinStatu
		isNew := true
		for _, coinstatu := range config.RealTimeStatus.CoinStatus {
			if coinstatu.Category == operate.CoinCategory {
				isNew = false
				coinstatu.Used = operate.CoinUsed
			}
			newCoinStatus = append(newCoinStatus, coinstatu)
		}

		if isNew {
			//TODO
			newCoinStatus = append(newCoinStatus, config.CoinStatu{Name: config.COIN_NAME_BTC, Category: operate.CoinCategory, Decimals: config.COIN_DECIMALS_BTC, Used: operate.CoinUsed})
		}
		reportInfo := &config.GrpcStream{Type: config.GRPC_COIN_LIST_WEB} //上报
		config.ReportedChan <- reportInfo

		config.RealTimeStatus.CoinStatus = newCoinStatus
		landStatus(handler.db, "")
	}

}

//key generate
func generateKey(db localdb.Database, reqType string, appNum int, appId string, pass string, code string, opinion bool) (bool, int) {

	defer apiMu.Unlock()
	apiMu.Lock()

	batchCount := batchCountMap[reqType]

	if batchCount == appNum { //上次输入完整
		batchCount = 0
		delete(batchCountMap, reqType)
		for k, _ := range authorizedMap { //删除授权信息
			if strings.HasPrefix(k, reqType) {
				delete(authorizedMap, k)
			}
		}
		for k, _ := range passMap { //删除
			if strings.HasPrefix(k, reqType) {
				delete(passMap, k)
			}
		}
	}

	batchCount++
	batchCountMap[reqType] = batchCount

	passMap[reqType+appId] = pass
	if reqType == config.REQ_CREATE { //创建
		codeMap[reqType+appId] = code
	}

	authorizedMap[reqType+appId] = opinion

	if batchCount < appNum { //未输入完整密码

	} else if batchCount == appNum { //全部输入完成
		if reqType == config.REQ_HASH { //hash确认
			for k, b := range authorizedMap {
				if strings.HasPrefix(k, reqType) {
					if !b { //有一个拒绝即失效
						return false, config.PASSWORD_STATUS_OK
					}
				}
			}
		}

		var passBytes [][]byte = make([][]byte, appNum)
		count := 0
		for k, v := range passMap {
			if strings.HasPrefix(k, reqType) {
				passBytes[count] = []byte(v)
				count++
			}
		}

		if reqType == config.REQ_CREATE { //创建
			firstCode := ""
			for _, code := range codeMap {
				if firstCode == "" {
					firstCode = code
				} else if firstCode != code {
					for k, _ := range codeMap { //删除
						if strings.HasPrefix(k, reqType) {
							delete(codeMap, k)
						}
					}
					return false, config.PASSWORD_CODE_FAILED
				}
			}

			if !trans.ExistPrivateKeyHash(db) { //不存在
				if trans.RecoverPrivateKey(db, passBytes...) == nil {
					return true, config.PASSWORD_STATUS_FAILED
				}
			}
		} else { //校验
			if trans.CheckPrivateKey(db, passBytes...) {
				return true, config.PASSWORD_STATUS_OK
			} else {
				return false, config.PASSWORD_STATUS_FAILED
			}
		}
	}
	return false, config.PASSWORD_STATUS_FAILED
}

//校验key
func checkKey(db localdb.Database, reqType string, appNum int, appId string, pass string, opinion bool) (bool, int) {
	return generateKey(db, reqType, appNum, appId, pass, "", opinion)
}

//状态持久化
func landStatus(db localdb.Database, reqType string) {
	if reqType != "" {
		var nodesAuthorized []config.NodeAuthorized //从新
		for k, v := range authorizedMap {
			if strings.HasPrefix(k, reqType) {
				nodesAuthorized = append(nodesAuthorized, config.NodeAuthorized{ApplyerId: strings.Trim(k, reqType), Authorized: v})
			}
		}
		config.RealTimeStatus.NodesAuthorized = nodesAuthorized
	}

	if statusByte, err := json.Marshal(config.RealTimeStatus); err != nil {
		log.Error("disallowHandler json marshal error:%v", err)
	} else {
		if err = db.Put([]byte(config.STATAUS_KEY), statusByte); err != nil {
			log.Error("landStatus err :%s", err)
		}
	}
}

func getKey(data, sign string, db localdb.Database, appId string) string {
	if err := util.RsaSignVer([]byte(data), []byte(sign), db, appId); err != nil {
		log.Error("get key err: %s", err)
	} else {
		return string(util.CBCDecrypter([]byte(data), util.GetAesKeyRandomFromDb(db)))
	}
	return ""
}

func checkSign(flow, sign, appId string, db localdb.Database) bool {
	if err := util.RsaSignVer([]byte(flow), []byte(sign), db, appId); err != nil {
		log.Info("check sign failed.")
		return false
	}
	return true
}
