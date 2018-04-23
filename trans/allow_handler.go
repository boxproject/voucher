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
	logger "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
)

func newAllowHandler(ctx *HandlerContext) EventHandler {
	return &allowHandler{ctx, true, allowEvent}
}

type allowHandler struct {
	context *HandlerContext
	flag    bool
	name    common.Hash
}

func (h *allowHandler) Name() common.Hash {
	return h.name
}

func (h *allowHandler) Scan(log *common.EtherLog) error {
	logger.Debug("allowHandler.......")
	if !common.AddressEquals(log.Address, h.context.ContractAddress) {
		return nil
	}

	// TODO 直接拿到data数据就是公链成功的签名流模版哈希
	//log.Data
	//log.BlockNumber
	//log.BlockHash
	//log.WDHash
	//_, err := h.context.State.Load(nil)
	//if err != nil {
	//	if err == errors.NoDataErr {
	//		logger.Warnf("the hash: %s not in local database.", common.BytesToHex(log.Data, true))
	//		return nil
	//	}
	//
	//	return err
	//}

	if dataBytes := log.Data; len(dataBytes) > 0 {
		hash := common.BytesToHash(dataBytes[:32])
		logger.Debug("allowHandler hash: %s", hash.Hex())
		if _, err := h.context.Db.Get([]byte(config.HASH_LIST_PRIFIX + hash.Hex())); err != nil {
			logger.Error("allowHandler load content err:%v", err)
		} else {
			//供查询使用
			//hashModel := &config.GrpcStream{}
			//if err := json.Unmarshal(hashByte, hashModel); err != nil {
			//	logger.Error("allowHandler json marshal error:%v", err)
			//	return err
			//}
			//hashModel.Status = config.HASH_STATUS_ENABLE
			config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_HASH_ENABLE_WEB, Hash: hash}
			//if newHashByte, err := json.Marshal(hashModel); err != nil {
			//	logger.Error("allowHandler json marshal error:%v", err)
			//} else {
			//	//h.context.State.Save([]byte(config.HASH_LIST_PRIFIX+hash.Hex()), newHashByte)
			//	hashHandler.HashList()
			//}
		}
	}
	return nil
}
