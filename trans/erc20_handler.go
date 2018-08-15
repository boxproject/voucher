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
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/token"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

func newERC20Handler(cfg *HandlerContext, cli *ethclient.Client) EventHandler {
	return &erc20Handler{cfg, cli}
}

type erc20Handler struct {
	context *HandlerContext
	client  *ethclient.Client
}

func (e *erc20Handler) Name() common.Hash {
	return erc20Event
}

func (e *erc20Handler) Scan(eLog *common.EtherLog) error {
	length := len(eLog.Topics)
	//log.Debug("~~~~~ ERC20 event ~~~~~ topic length: %d", length)
	if length == 3 {

		// May be ERC20
		erc20Addr := eLog.Address.Hex()
		// 检测erc20addr 是否是合法的需要监控的地址
		isErc20, category := scanAddress(erc20Addr)
		if !isErc20 {
			return nil
		}
		// erc20 存在并被监管

		// 检查是否充值成功
		// to 是否是监控的交易所的合约地址
		// 如果是，则需要关联to 对应的 uuid
		// 并登记一笔充值成功的记录
		from := common.BytesToAddress(eLog.Topics[1].Bytes()[:32]).Hex()
		to := common.BytesToAddress(eLog.Topics[2].Bytes()[:32]).Hex()

		// amount
		amount := common.BytesToHash(eLog.Data[:32]).Big()

		log.Info("erc20Handler category: %d, from: %s, to: %s, amount: %d", category, from, to, amount)

		//if wallet, err := NewWallet(common.HexToAddress(to), e.client); err != nil {
		//	log.Error("wallet error: %s", err)
		//} else {
		//	opts := NewKeyedTransactor()
		//	wallet.TransferERC20(opts, eLog.Address, e.context.ContractAddress, amount)
		//}

		contractAddr := config.RealTimeStatus.ContractAddress

		if from == contractAddr { //提现
			if wdHash, err := e.context.Db.Get([]byte(config.WITHDRAW_TX_PRIFIX + eLog.TxHash.Hex())); err != nil {
				log.Error("get db err: %s", err)
			} else {
				config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_WITHDRAW_WEB, From: from, To: to, Amount: amount, TxHash: eLog.TxHash.Hex(), WdHash: common.HexToHash(string(wdHash)), Category: big.NewInt(category)}
			}
		} else if to == contractAddr { //充值
			config.ReportedChan <- &config.GrpcStream{Type: config.GRPC_DEPOSIT_WEB, From: from, To: to, Amount: amount, TxHash: eLog.TxHash.Hex(), Category: big.NewInt(category)}
		}
	}
	return nil
}

func scanAddress(addr string) (bool, int64) {
	if ethToken := token.GetTokenByAddr(strings.ToLower(addr)); ethToken != nil {
		return true, ethToken.Category
	}
	return false, 0
}
