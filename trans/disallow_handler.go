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
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	logger "github.com/alecthomas/log4go"
)

//func newDisallowHandler(context *HandlerContext) EventHandler {
//	return &disallowHandler{context}
//}

type disallowHandler struct {
	context *HandlerContext
	flag    bool
	name    common.Hash
}

func (handler *disallowHandler) Name() common.Hash {
	return disallowEvent
}

func newDisallowHandler(cfg *HandlerContext) EventHandler {
	return &disallowHandler{cfg, false, disallowEvent}
}

func (h *disallowHandler) Scan(log *common.EtherLog) error {
	logger.Debug("disallowHandler.......")
	if !common.AddressEquals(log.Address, h.context.ContractAddress) {
		return nil
	}

	if dataBytes := log.Data; len(dataBytes) > 0 {
		hash := common.BytesToHash(dataBytes[:32])
		logger.Debug("disallowHandler, hash: %s", hash)
		if _, err := h.context.Db.Get([]byte(config.HASH_LIST_PRIFIX + hash.Hex())); err != nil {
			logger.Error("disallowHandler load content err:%v", err)
		} else {
			//供查询使用
			//hashModel := &config.GrpcStream{}
			//if err := json.Unmarshal(hashByte, hashModel); err != nil {
			//	logger.Error("disallowHandler json marshal error:%v", err)
			//	return err
			//}
			config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_HASH_ENABLE_WEB, Hash: hash}
			//hashModel.Status = config.HASH_STATUS_DISABLE
			//if newHashByte, err := json.Marshal(hashModel); err != nil {
			//	logger.Error("disallowHandler json marshal error:%v", err)
			//} else {
			//	h.context.State.Save([]byte(config.HASH_LIST_PRIFIX+hash), newHashByte)
			//	hashHandler.HashList(h.context.Db)
			//}

		}
	}
	return nil
}
