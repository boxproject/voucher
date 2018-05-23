package token

import (
	"encoding/json"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	"sync"
)

var tokenRWMutex sync.RWMutex

var EthAddrTokenMap map[string]*config.TokenInfo = make(map[string]*config.TokenInfo) //eth addr-token map

var EthCategoryTokenMap map[int64]*config.TokenInfo = make(map[int64]*config.TokenInfo) //eth category-token map

//新增或编辑代币
func AddTokenMap(tokenInfoD *config.TokenInfo, db localdb.Database) (bool, int) {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()

	log.Debug("EthAddrTokenMap...", EthAddrTokenMap)

	tokenInfoS := EthAddrTokenMap[tokenInfoD.ContractAddr]
	if tokenInfoS != nil { //存在
		tokenInfoD.Category = tokenInfoS.Category
	} else {

		tokenLen := len(EthAddrTokenMap)
		if tokenLen == 0 {
			tokenLen = 2
		} else {
			tokenLen = tokenLen + 1 + 1 //TODO 代币从2开始
		}
		tokenInfoD.Category = int64(tokenLen)
	}

	if tokenData, err := json.Marshal(tokenInfoD); err != nil {
		log.Error("token marshal err: %s", err)
	} else {
		if err = db.Put([]byte(config.TOKEN_PRIFIX+tokenInfoD.ContractAddr), tokenData); err != nil {
			log.Error("token db land err: %s", err)
		} else {
			EthAddrTokenMap[tokenInfoD.ContractAddr] = tokenInfoD
			EthCategoryTokenMap[tokenInfoD.Category] = tokenInfoD
			//TODO 上报
			reportInfo := &config.GrpcStream{Type: config.GRPC_TOKEN_LIST_WEB}
			if tokenMap, err := db.GetPrifix([]byte(config.TOKEN_PRIFIX)); err != nil {
				log.Error("get tokenlist err:%s", err)
			} else {
				for _, tokenBytes := range tokenMap {
					tokenInfo := &config.TokenInfo{}
					if err = json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
						log.Error("unmarshal err:%s", err)
					} else {
						reportInfo.TokenList = append(reportInfo.TokenList, tokenInfo)
					}
				}
				config.ReportedChan <- reportInfo
			}
			return true, len(EthAddrTokenMap)
		}
	}
	return false, len(EthAddrTokenMap)
}

func DelTokenMap(tokenInfoD *config.TokenInfo, db localdb.Database) (bool, int) {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()

	tokenInfoS := EthAddrTokenMap[tokenInfoD.ContractAddr]
	if tokenInfoS != nil { //存在
		db.Delete([]byte(config.TOKEN_PRIFIX + tokenInfoD.ContractAddr))
		delete(EthAddrTokenMap, tokenInfoS.ContractAddr)
		delete(EthCategoryTokenMap, tokenInfoS.Category)
		//TODO 上报
		reportInfo := &config.GrpcStream{Type: config.GRPC_TOKEN_LIST_WEB}
		if tokenMap, err := db.GetPrifix([]byte(config.TOKEN_PRIFIX)); err != nil {
			log.Error("get tokenlist err:%s", err)
		} else {
			for _, tokenBytes := range tokenMap {
				tokenInfo := &config.TokenInfo{}
				if err = json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
					log.Error("unmarshal err:%s", err)
				} else {
					reportInfo.TokenList = append(reportInfo.TokenList, tokenInfo)
				}
			}
			config.ReportedChan <- reportInfo
		}
		return true, len(EthAddrTokenMap)
	} else {
		log.Error("del token failed. not found token address: %s", tokenInfoD.ContractAddr)
	}
	return false, len(EthAddrTokenMap)
}

func TokenList(db localdb.Database) bool {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()

	reportInfo := &config.GrpcStream{Type: config.GRPC_TOKEN_LIST_WEB}
	if tokenMap, err := db.GetPrifix([]byte(config.TOKEN_PRIFIX)); err != nil {
		log.Error("get tokenlist err:%s", err)
		return false
	} else {
		for _, tokenBytes := range tokenMap {
			tokenInfo := &config.TokenInfo{}
			if err = json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
				log.Error("unmarshal err:%s", err)
			} else {
				reportInfo.TokenList = append(reportInfo.TokenList, tokenInfo)
			}
		}
		config.RealTimeStatus.TokenCount = len(tokenMap)
		if statusByte, err := json.Marshal(config.RealTimeStatus); err != nil {
			log.Error("hashList json marshal error:%v", err)
		} else {
			if err = db.Put([]byte(config.STATAUS_KEY), statusByte); err != nil {
				log.Error("landStatus err :%s", err)
			}
		}
		config.ReportedChan <- reportInfo
	}
	return true

}

//db加载token信息
func LoadTokenFrom(db localdb.Database) {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()
	log.Debug("LoadTokenFrom....", EthAddrTokenMap)
	if tokenMap, err := db.GetPrifix([]byte(config.TOKEN_PRIFIX)); err != nil {
		log.Error("get tokenlist err:%s", err)
	} else {
		for _, tokenBytes := range tokenMap {
			tokenInfo := &config.TokenInfo{}
			if err = json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
				log.Error("LoadTokenFrom unmarshal err: %s", err)
			} else {
				EthAddrTokenMap[tokenInfo.ContractAddr] = tokenInfo
				EthCategoryTokenMap[tokenInfo.Category] = tokenInfo
				log.Error("unmarshal err:%s", err)
			}
		}
	}

	log.Debug("EthAddrTokenMap....", EthAddrTokenMap)
}

func GetTokenByAddr(addr string) *config.TokenInfo {
	tokenRWMutex.Lock()
	log.Debug("EthAddrTokenMap....", EthAddrTokenMap)
	defer tokenRWMutex.Unlock()
	return EthAddrTokenMap[addr]
}

//根据category获取token address
func GetTokenByCategory(category int64) *config.TokenInfo {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()
	return EthCategoryTokenMap[category]
}
