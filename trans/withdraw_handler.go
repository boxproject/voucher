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

func newWithdrawHandler(context *HandlerContext) EventHandler {
	return &withdrawHandler{context}
}

type withdrawHandler struct {
	context *HandlerContext
}

func (h *withdrawHandler) Name() common.Hash {
	return withdrawEvent
}

func (h *withdrawHandler) Scan(log *common.EtherLog) error {
	logger.Debug("withdrawHandler scan...")
	if !common.AddressEquals(log.Address, h.context.ContractAddress) {
		return nil
	}

	if dataBytes := log.Data; len(dataBytes) > 0 {
		address := common.BytesToAddress(dataBytes[:32]).Hex()
		amount := common.BytesToHash(dataBytes[32:64]).Big()
		if wdHash, err := h.context.Db.Get([]byte(config.WITHDRAW_TX_PRIFIX + log.TxHash.Hex())); err != nil {
			logger.Error("get db err: %s", err)
		} else {
			config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_WITHDRAW_WEB,  To: address, Amount: amount, WdHash: common.HexToHash(string(wdHash)),TxHash: log.TxHash.Hex()}
		}
	}
	return nil
}
