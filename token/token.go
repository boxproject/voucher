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

//获取有效token
func getEnableTokenMap(db localdb.Database) ([]*config.TokenInfo){
	var TokenList []*config.TokenInfo
	if tokenMap, err := db.GetPrifix([]byte(config.TOKEN_PRIFIX)); err != nil {
		log.Error("get tokenlist err:%s", err)
	} else {
		for _, tokenBytes := range tokenMap {
			tokenInfo := &config.TokenInfo{}
			if err = json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
				log.Error("unmarshal err:%s", err)
			} else {
				if tokenInfo.Status == "true" {
					TokenList = append(TokenList, tokenInfo)
				}
			}
		}
	}
	return TokenList
}
//
func getEthTokenCount(db localdb.Database)(allCount,enbleCount,disableCount int){
	if tokenMap, err := db.GetPrifix([]byte(config.TOKEN_PRIFIX)); err != nil {
		log.Error("get tokenlist err:%s", err)
	} else {
		//allCount = len(tokenMap)
		for _, tokenBytes := range tokenMap {
			tokenInfo := &config.TokenInfo{}
			if err = json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
				log.Error("unmarshal err:%s", err)
			} else {
				if tokenInfo.Status == "true" {
					enbleCount++
				}else {
					disableCount++
				}
			}
		}
	}
	return disableCount+enbleCount,enbleCount,disableCount
}


//新增或编辑代币
func AddTokenMap(tokenInfoD *config.TokenInfo, db localdb.Database) (bool, int) {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()

	log.Debug("EthAddrTokenMap...", EthAddrTokenMap)
	tokenInfoD.Status = "true"
	tokenInfo := &config.TokenInfo{}
	tokenBytes,_ := db.Get([]byte(config.TOKEN_PRIFIX+tokenInfoD.ContractAddr))
	if len(tokenBytes) > 0 {
		if err := json.Unmarshal([]byte(tokenBytes), tokenInfo); err != nil {
			log.Error("unmarshal err:%s", err)
			return false, len(EthAddrTokenMap)
		} else {
			log.Debug("get tokenInfo:",tokenInfo)
			tokenInfoD.Category = tokenInfo.Category
		}
	}else {
		//add new token
		tokenCount,_,_ := getEthTokenCount(db)
		tokenCount = tokenCount + 2 //start from 2
		tokenInfoD.Category = int64(tokenCount)
	}
	tokenData, err := json.Marshal(tokenInfoD)
	if err != nil {
		log.Error("token marshal err: %s", err)
		return false, len(EthAddrTokenMap)
	}
	if err = db.Put([]byte(config.TOKEN_PRIFIX+tokenInfoD.ContractAddr), tokenData); err != nil {
		log.Error("token db land err: %s", err)
	} else {
		EthAddrTokenMap[tokenInfoD.ContractAddr] = tokenInfoD
		EthCategoryTokenMap[tokenInfoD.Category] = tokenInfoD
		//TODO 上报
		reportInfo := &config.GrpcStream{Type: config.GRPC_TOKEN_LIST_WEB}
		reportInfo.TokenList = getEnableTokenMap(db)
		config.ReportedChan <- reportInfo

		return true, len(reportInfo.TokenList)
	}
	_,tokenEnableCount,_ := getEthTokenCount(db)
	return false, tokenEnableCount
}

func DelTokenMap(tokenInfoD *config.TokenInfo, db localdb.Database) (bool, int) {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()

	tokenInfoS := EthAddrTokenMap[tokenInfoD.ContractAddr]
	if tokenInfoS != nil { //存在
		//db.Delete([]byte(config.TOKEN_PRIFIX + tokenInfoD.ContractAddr))
		tokenInfoS.Status = "false"
		tokenData, err := json.Marshal(tokenInfoS)
		if  err != nil {
			log.Error("token marshal err: %s", err)
		}
		db.Put([]byte(config.TOKEN_PRIFIX+tokenInfoS.ContractAddr), tokenData)

		delete(EthAddrTokenMap, tokenInfoS.ContractAddr)
		delete(EthCategoryTokenMap, tokenInfoS.Category)
		//TODO 上报
		reportInfo := &config.GrpcStream{Type: config.GRPC_TOKEN_LIST_WEB}
		reportInfo.TokenList = getEnableTokenMap(db)
		config.ReportedChan <- reportInfo

		return true, len(reportInfo.TokenList)
	} else {
		log.Error("del token failed. not found token address: %s", tokenInfoD.ContractAddr)
	}

	_,tokenEnableCount,_ := getEthTokenCount(db)
	return false, tokenEnableCount
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
				if tokenInfo.Status == "true" {
					reportInfo.TokenList = append(reportInfo.TokenList, tokenInfo)
				}
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
				//db.Delete([]byte(config.TOKEN_PRIFIX+tokenInfo.ContractAddr))
				//continue
				if tokenInfo.Status == "true"  {
					EthAddrTokenMap[tokenInfo.ContractAddr] = tokenInfo
					EthCategoryTokenMap[tokenInfo.Category] = tokenInfo
				}else if len(tokenInfo.Status) == 0{
					//DB兼容处理
					log.Info("[OLD TOKEN DB][%v]",tokenInfo)
					tokenInfo.Status = "true"
					EthAddrTokenMap[tokenInfo.ContractAddr] = tokenInfo
					EthCategoryTokenMap[tokenInfo.Category] = tokenInfo
					log.Info("[DO UPDATE][%v]",tokenInfo)
					tokenData, err := json.Marshal(tokenInfo)
					if err != nil {
						log.Error("token marshal err: %s", err)
					}
					if err = db.Put([]byte(config.TOKEN_PRIFIX+tokenInfo.ContractAddr), tokenData); err != nil {
						log.Error("token db land err: %s", err)
					}
				}
			}
		}
	}

	log.Debug("EthAddrTokenMap....", EthAddrTokenMap)
}

func GetTokenByAddr(addr string) *config.TokenInfo {
	tokenRWMutex.Lock()
	//log.Debug("EthAddrTokenMap....", EthAddrTokenMap)
	defer tokenRWMutex.Unlock()
	return EthAddrTokenMap[addr]
}

//根据category获取token address
func GetTokenByCategory(category int64) *config.TokenInfo {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()
	return EthCategoryTokenMap[category]
}
