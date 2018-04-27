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

package node

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	"encoding/json"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	"github.com/boxproject/voucher/pb"
	"github.com/boxproject/voucher/util"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//对外api
type NodeApi struct {
	cfg         config.AgentServiceConfig
	db          localdb.Database
	conn        *grpc.ClientConn
	heartQuitCh chan struct{}
	quitCh      chan struct{}
}

type HashResult struct {
	NodeName   string
	NodeResult bool
}

//init
func InitConn(cfg *config.Config, db localdb.Database) (*NodeApi, error) {
	log.Info("InitConn....")
	nodeApi := &NodeApi{cfg: cfg.AgentSerCfg, db: db, heartQuitCh: make(chan struct{}, 1), quitCh: make(chan struct{}, 1)}
	cred, err := nodeApi.loadCredential(cfg.AgentSerCfg.ClientCertPath, cfg.AgentSerCfg.ClientKeyPath)
	if err != nil {
		log.Error("load tls cert failed. cause: %v\n", err)
		return nil, err
	}
	nodeApi.conn, err = grpc.Dial(cfg.AgentSerCfg.RPCAddr, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Error("connect to the remote server failed. cause: %v", err)
		return nil, err
	}
	go nodeApi.streamRecv()
	go nodeApi.RepStart() //启动接受充值上报请求
	return nodeApi, nil
}

//close conn
func (n *NodeApi) CloseConn() error {
	log.Info("conn closed begin")
	if n.conn == nil {
		return errors.New("conn nil")
	} else {
		if n.quitCh != nil { //流服务
			n.quitCh <- struct{}{}
		}
		if n.heartQuitCh != nil { //心跳
			n.heartQuitCh <- struct{}{}
		}
		n.conn.Close()
	}
	log.Info("conn closed end")
	return nil
}

//report
func (n *NodeApi) RepStart() {
	log.Info("node rep started!")
	for true {
		select {
		case data, ok := <-config.ReportedChan:
			if ok {
				log.Debug("data.Type....", data.Type)
				switch data.Type {
				case config.GRPC_DEPOSIT_WEB:
					n.DepositRequest(data)
					break
				case config.GRPC_WITHDRAW_TX_WEB:
					n.WithDrewTxRequest(data)
					break
				case config.GRPC_WITHDRAW_WEB:
					n.WithDrewRequest(data)
					break
				//case config.GRPC_HASH_LIST_WEB:
				//	n.HashListRequest(data)
				//	break
				case config.GRPC_TOKEN_LIST_WEB:
					n.TokenListRequest(data)
					break
				case config.GRPC_COIN_LIST_WEB:
					n.CoinListRequest(data)
					break
				case config.GRPC_HASH_ENABLE_WEB:
					n.HashEnableRequest(data)
					break
				default:
					log.Info("unknown node req type:%v", data.Type)
				}
			} else {
				log.Error("read node from channel failed")
			}
		}
	}
}

func (n *NodeApi) heart() {
	log.Info("nodeApi heart started!")
	ticker := time.NewTicker(5 * time.Second)
LOOP:
	for {
		select {
		case <-ticker.C:
			client := pb.NewSynchronizerClient(n.conn)
			if statuBytes, err := json.Marshal(config.RealTimeStatus); err != nil {
				log.Error("json marshal err: %s", err)
			} else {
				if rsp, err := client.Heart(context.TODO(), &pb.HeartRequest{RouterType: "grpc", ServerName: n.cfg.Name, Name: n.cfg.Alias, Ip: util.GetCurrentIp(), Msg: statuBytes}); err != nil {
					log.Error("heart req failed %s\n", err)
				} else {
					log.Debug("heart response", rsp)
				}
			}

		case <-n.heartQuitCh:
			ticker.Stop()
			break LOOP
		}
	}

	log.Info("nodeApi heart stopped!")

}

////同意
//func (n *NodeApi) AllowRequest(msg []byte) bool {
//	n.router(config.ROUTER_TYPE_GRPC, n.cfg.CompanionName, msg)
//	return true
//}
//
////禁用
//func (n *NodeApi) DisAllowRequest(msg []byte) bool {
//	n.router(config.ROUTER_TYPE_GRPC, n.cfg.CompanionName, msg)
//	return true
//}

//充值上报
func (n *NodeApi) DepositRequest(grpcStream *config.GrpcStream) bool {
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("DepositRequest marshal err: %s", err)
		//TODO
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

//提现tx上报
func (n *NodeApi) WithDrewTxRequest(grpcStream *config.GrpcStream) bool {
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("WithDrewTxRequest marshal err: %s", err)
		//TODO
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

//提现上报
func (n *NodeApi) WithDrewRequest(grpcStream *config.GrpcStream) bool {
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("WithDrewRequest marshal err: %s", err)
		//TODO
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

//hash上报
func (n *NodeApi) HashListRequest(grpcStream *config.GrpcStream) bool {
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("HashListRequest marshal err: %s", err)
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

//token上报
func (n *NodeApi) TokenListRequest(grpcStream *config.GrpcStream) bool {
	log.Debug("TokenListRequest......")
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("TokenListRequest marshal err: %s", err)
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

//list上报
func (n *NodeApi) CoinListRequest(grpcStream *config.GrpcStream) bool {
	log.Debug("CoinListRequest......")
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("CoinListRequest marshal err: %s", err)
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

//hash 确认上报
func (n *NodeApi) HashEnableRequest(grpcStream *config.GrpcStream) bool {
	log.Debug("HashEnableRequest......")
	if msg, err := json.Marshal(grpcStream); err != nil {
		log.Error("HashEnableRequest marshal err: %s", err)
	} else {
		n.router(config.ROUTER_TYPE_WEB, "", msg)
	}
	return true
}

func (n *NodeApi) router(routerType string, routerName string, msg []byte) (string, error) {
	log.Debug("router....")
	if n.conn != nil {
		client := pb.NewSynchronizerClient(n.conn)
		if rsp, err := client.Router(context.TODO(), &pb.RouterRequest{routerType, routerName, msg}); err != nil {
			log.Error("depositTx req failed %v\n", err)
			return "", err
		} else {
			return rsp.Code, nil
		}
	}
	return "", errors.New("conn is nil")
}

//stream recv
func (n *NodeApi) streamRecv() {
	timeCount := 1
	for { //5秒尝试连接一次 TODO
		log.Info("try...", timeCount)
		client := pb.NewSynchronizerClient(n.conn)
		stream, err := client.Listen(context.TODO())
		if err != nil { //5秒尝试连接一次
			log.Error("[STREAM ERR] %v\n", err)
		} else {
			//获取上次block number
			waitc := make(chan struct{})
			stream.Send(&pb.ListenReq{ServerName: n.cfg.Name, Name: n.cfg.Alias, Ip: util.GetCurrentIp()})
			go func() {
				for {
					if resp, err := stream.Recv(); err != nil { //rec error
						log.Error("[STREAM ERR] %v\n", err)
						close(waitc)
						return
					} else {
						log.Debug("stream Recv: %s\n", resp)
						n.handleStream(resp)
					}
				}
			}()
			go n.heart() //启动心跳
			<-waitc
			if err = stream.CloseSend(); err != nil {
				log.Error("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
			}
		}
		time.Sleep(5 * time.Second)
		timeCount++
	}
	log.Debug("streamRecv stop")
}

//cert
func (n *NodeApi) loadCredential(clientCertPath, clientKeyPath string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(clientCertPath)
	if err != nil {
		return nil, err
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, err
	}

	config := &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	return credentials.NewTLS(config), nil
}

//处理流
func (n *NodeApi) handleStream(streamRsp *pb.StreamRsp) {
	streamModel := &config.GrpcStream{}
	if err := json.Unmarshal(streamRsp.Msg, streamModel); err != nil {
		log.Error("json marshal error:%v", err)
		return
	}

	switch streamModel.Type {
	//case config.GRPC_HASH_ADD_REQ: //新加hash
	//	n.db.Put([]byte(config.HASH_LIST_PRIFIX+streamModel.Hash.Hex()), streamRsp.Msg) //供查询使用
	//hash.HashList(n.db)
	//case config.GRPC_HASH_ENABLE_REQ: //确认hash req
	//	n.AllowRequest(streamRsp.Msg)
	//case config.GRPC_HASH_DISABLE_REQ: //拒绝hash req
	//	n.DisAllowRequest(streamRsp.Msg)
	case config.GRPC_HASH_ADD_LOG: //hash log
		n.db.Put([]byte(config.HASH_LIST_PRIFIX+streamModel.Hash.Hex()), []byte(streamModel.Flow)) //供查询使用
	case config.GRPC_HASH_ENABLE_LOG: //同意 log
		if checkSigns(streamModel.Flow, streamModel.SignInfos, n.db) {
			n.db.Put([]byte(config.HASH_LIST_PRIFIX+streamModel.Hash.Hex()), []byte(streamModel.Flow))           //供查询使用
			n.db.Put([]byte(config.PENDING_PRIDFIX+streamModel.Hash.Hex()), streamRsp.Msg)                       //pending
			config.Ecr20RecordChan <- &config.Ecr20Record{Type: config.ECR20_TYPE_ALLOW, Hash: streamModel.Hash} //上公链
		}
	case config.GRPC_HASH_DISABLE_LOG: //禁用 log
		if checkSigns(streamModel.Flow, streamModel.SignInfos, n.db) { //pending
			n.db.Put([]byte(config.PENDING_PRIDFIX+streamModel.Hash.Hex()), streamRsp.Msg)
			config.Ecr20RecordChan <- &config.Ecr20Record{Type: config.ECR20_TYPE_DISALLOW, Hash: streamModel.Hash} //上公链
		}
	case config.GRPC_WITHDRAW_LOG: //提现 log
		//TODO 验签
		if checkWithdrawSign(streamModel.Hash.Hex(), streamModel.Flow, streamModel.WdFlow, n.db) {
			log.Debug("check sign success！")
			n.db.Put([]byte(config.WITHDRAW_APPLY_PRIFIX+streamModel.WdHash.Hex()), streamRsp.Msg) // for view
			n.db.Put([]byte(config.PENDING_PRIDFIX+streamModel.WdHash.Hex()), streamRsp.Msg)       //pending
			category := streamModel.Category
			if category != nil {
				switch streamModel.Category.Int64() {
				case config.CATEGORY_BTC:
					//bitcoin
					config.BtcRecordChan <- &config.BtcRecord{Type: config.BTC_TYPE_APPROVE, WdHash: streamModel.WdHash, Handlefee: streamModel.Fee, To: streamModel.To, Amount: streamModel.Amount}
					break
				case config.CATEGORY_ETH:
					//上公链
					config.Ecr20RecordChan <- &config.Ecr20Record{Type: config.ECR20_TYPE_APPROVE, Hash: streamModel.Hash, WdHash: streamModel.WdHash, To: common.HexToAddress(streamModel.To), Amount: streamModel.Amount, Category: streamModel.Category}
					break
				default: //ETH 代币
					config.Ecr20RecordChan <- &config.Ecr20Record{Type: config.ECR20_TYPE_APPROVE, Hash: streamModel.Hash, WdHash: streamModel.WdHash, To: common.HexToAddress(streamModel.To), Amount: streamModel.Amount, Category: streamModel.Category}
				}
			}
		}
	case config.GRPC_VOUCHER_OPR_REQ: //签名机操作
		config.OperateChan <- streamModel.VoucherOperate
	default:

	}
}

//验证签名
func checkSigns(flow string, signInfos []*config.SignInfo, db localdb.Database) bool {
	var appMap map[string]interface{} = make(map[string]interface{})
	if len(signInfos) == config.RealTimeStatus.Total {
		for _, signInfo := range signInfos {
			appMap[signInfo.AppId] = nil
			if !checkSign(flow, signInfo.Sign, signInfo.AppId, db) {
				return false
			}
		}
		if len(appMap) == config.RealTimeStatus.Total { //私钥app都填完方通过
			log.Info("sign passed！")
			return true
		}
	} else {
		log.Error("sign failed！")
	}
	return false
}

func checkSign(flow, sign, appId string, db localdb.Database) bool {
	if err := util.RsaSignVer([]byte(flow), []byte(sign), db, appId); err != nil {
		log.Info("check sign failed.")
		return false
	}
	return true
}

//提现验签
func checkWithdrawSign(hashStr string, wdFlow, wdSign string, db localdb.Database) bool {
	var pubKeyMap map[string]string = make(map[string]string)
	if flowMsg, err := db.Get([]byte(config.HASH_LIST_PRIFIX + hashStr)); err != nil { //查询hash模板
		log.Error("db get err: %s", err)
	} else {
		hashFlow := &config.HashFlow{}
		if err := json.Unmarshal(flowMsg, hashFlow); err != nil {
			log.Error("json unmarshal error:%v", err)
			return false
		} else {
			for _, approvalInfo := range hashFlow.Approval_info {
				for _, approver := range approvalInfo.Approvers {
					pubKeyMap[approver.App_account_id] = approver.Pub_key
				}
			}
		}
	}

	var employeeFlows []config.EmployeeFlow
	if err := json.Unmarshal([]byte(wdSign), &employeeFlows); err != nil {
		log.Error("json unmarshal error:%v", err)
		return false
	} else {
		if employeeFlows == nil {
			log.Error("sign data nil！")
			return false
		}
		for _, employeeFlow := range employeeFlows {
			if pubKeyMap[employeeFlow.Appid] == "" {
				return false
			} else {
				if err = util.RsaSignVerWithKey([]byte(wdFlow), []byte(employeeFlow.Sign), []byte(pubKeyMap[employeeFlow.Appid])); err != nil {
					log.Error("sign check err: %s. [appId：%s]", err, employeeFlow.Appid)
					return false
				}
			}
		}
	}
	return true
}
