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
	logger "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func newDepositHandler(ctx *HandlerContext, cli *ethclient.Client) EventHandler {
	return &depositHandler{ctx, cli}
}

type depositHandler struct {
	context *HandlerContext
	client  *ethclient.Client
}

func (h *depositHandler) Name() common.Hash {
	return depositEvent
}

func (h *depositHandler) Scan(log *common.EtherLog) error {
	logger.Debug("depositHandler scan...")
	if !common.AddressEquals(log.Address, h.context.ContractAddress) {
		return nil
	}

	var from common.Address
	var toTx string

	block, err := h.client.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
	if err != nil {
		logger.Error("blockByNumber error: %s", err)
	}

	for _, tx := range block.Transactions() {
		if tx.Hash().Hex() == log.TxHash.Hex() {
			signer := types.NewEIP155Signer(tx.ChainId())
			from, err = signer.Sender(tx)
			if err != nil {
				logger.Error("blockByNumber getFrom err: %s", err)
			}
			toTx = tx.To().Hex()
		}
	}

	if dataBytes := log.Data; len(dataBytes) > 0 {
		to := common.BytesToAddress(dataBytes[:32]).Hex()
		if from.Hex() == to { //直充
			to = toTx
		}
		amount := common.BytesToHash(dataBytes[32:64]).Big()
		logger.Debug("depositHandler from: %s, to: %s, amount: %d", from.Hex(), to, amount)
		config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_DEPOSIT_WEB, From: from.Hex(), To: to, Amount: amount, TxHash: log.TxHash.Hex(), Category: big.NewInt(config.CATEGORY_ETH)}
	}
	return nil
}
