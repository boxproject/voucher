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
package config

import (
	"bytes"

	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type (
	MetaKey   []byte
	InitStage []byte
)

var (
	// 初始状态
	INIT = MetaKey{0}
	// 服务器端的d
	SECRET = MetaKey{1}
	// 一次性密码,计数
	PASS = MetaKey{2}
	// 当前区块号
	BLKNUMBER = MetaKey{3}
	// 公链上的智能合约地址
	BANKADDRESS = MetaKey{5}

	PRIVATEKEYHASH = MetaKey{6}
	// bitcoin master pubkey addr
	BTCMASTERPUBADDRESS = MetaKey{7}

	Inited = InitStage{1}

	DefPubEthGasLimit = 3000000 //默认gaslimit
)

const (
	TYPE_CONTRACT = "contract"
	TYPE_SERVICE  = "service"
)

//pri db key
const (
	HASH_LIST_PRIFIX      = "v_hl_" //hash list
	WITHDRAW_APPLY_PRIFIX = "v_wa_" //withdraw apply
	PENDING_PRIDFIX       = "p_"    //pending
)

const (
	AES_KEY        = "aes_key"
	APP_KEY_PRIFIX = "app_key_"
	STATAUS_KEY    = "status_key"
	TOKEN_PRIFIX   = "token_key_"
)

//pub db key
const (
	PUBKEY_ETH         = "pub_key_eth"
	PUBKEY_BTC         = "pub_key_btc"
	DEPOSIT_PRIFIX     = "pub_dp_"
	WITHDRAW_PRIFIX    = "pub_wd_"
	ALLOWKP_PRIFIX     = "pub_al_"
	DISALLOW_PRIFIX    = "pub_da_"
	WITHDRAW_TX_PRIFIX = "pub_wd_tx_"
)

//确认、禁用请求 已作废
const (
	REQ_HASH = "hash"
	//REQ_HASH_DISALLOW = "disallow"
	REQ_CREATE = "create"
	REQ_DEPLOY = "deploy"
	REQ_START  = "start"
	REQ_PAUSE  = "pause"
)

const (
	Err_OK        = "0"  //正确
	Err_UNKONWN   = "1"  //未知的错误
	Err_MARSHAL   = "10" //josn打包失败
	Err_DB        = "11" //db失败
	Err_SERVER    = "12" //服务未启动
	Err_SERVER_ON = "13" //服务已启动
	Err_ADDR      = "14" //addr不存在
	Err_CONTRACT  = "15" //合约未发布

	Err_ETH             = "100" //链处理失败
	Err_UNENABLE_PREFIX = "101" //非法hash前缀
	Err_UNENABLE_LENGTH = "102" //非法hash值长度
	Err_KEY_RECOVER     = "103" //生成key失败
	Err_CODE_ROLE_NULL  = "104" //code为空
	Err_CODE_INCONS     = "105" //code不一致

	Err_BTC = "106" //BTC创建失败

)

const (
	LAST_BLOCK_PRIFX = "lbp_" //block前缀
)

const (
	CHAN_MAX_SIZE    = 100000 //chan 默认大小
	POOL_SIZE        = 2      //账户池大小
	NONCE_PLUS       = 1      //nonce值偏移值
)

const (
	ROUTER_TYPE_WEB  = "web"
	ROUTER_TYPE_GRPC = "grpc"
)

//grpc流类型
const (
	GRPC_HASH_ADD_REQ     = "1"  //hash add申请
	GRPC_HASH_ADD_LOG     = "2"  //hans add 私链log
	GRPC_HASH_ENABLE_REQ  = "3"  //hash enable 申请
	GRPC_HASH_ENABLE_LOG  = "4"  //hash enable 私链log
	GRPC_HASH_DISABLE_REQ = "5"  //hash disable 申请
	GRPC_HASH_DISABLE_LOG = "6"  //hash disable 私链log
	GRPC_WITHDRAW_REQ     = "7"  //提现 申请
	GRPC_WITHDRAW_LOG     = "8"  //提现 私链log
	GRPC_DEPOSIT_WEB      = "9"  //充值上报
	GRPC_WITHDRAW_TX_WEB  = "10" //提现tx上报
	GRPC_WITHDRAW_WEB     = "11" //	提现结果上报
	GRPC_VOUCHER_OPR_REQ  = "12" //	签名机操作处理
	//GRPC_HASH_LIST_REQ    = "13" //	审批流查询
	//GRPC_HASH_LIST_WEB    = "14" //	审批流上报
	GRPC_TOKEN_LIST_WEB   = "15" //	token上报
	GRPC_COIN_LIST_WEB    = "16" //coin上报
	GRPC_HASH_ENABLE_WEB  = "17" //hash enable 公链log
	GRPC_HASH_DISABLE_WEB = "18" //hash enable 公链log
)

//公链操作类型
const (
	ECR20_TYPE_ALLOW    = "1"
	ECR20_TYPE_DISALLOW = "2"
	ECR20_TYPE_APPROVE  = "3"
	BTC_TYPE_APPROVE    = ECR20_TYPE_APPROVE
)

//上报
const (
	REP_ACCOUNT_ADD = "1"
	REP_DEPOSIT     = "2"
	REP_WITHDRAW    = "3" //提现tx成功上报
	REP_WITHDRAW_TX = "4" //提现txid上报
)

//签名机状态
const (
	VOUCHER_STATUS_UNCONNETED = 0 //未连接
	VOUCHER_STATUS_UNCREATED  = 1 //未创建
	VOUCHER_STATUS_CREATED    = 2 //已创建
	VOUCHER_STATUS_DEPLOYED   = 3 //已发布
	VOUCHER_STATUS_STATED     = 4 //已启动
	VOUCHER_STATUS_PAUSED     = 5 //已停止
)

//关键句状态
const (
	PASSWORD_STATUS_OK     = 0 //成功
	PASSWORD_STATUS_FAILED = 1 //失败
	PASSWORD_SYSTEM_FAILED = 2 //系统异常
)

const (
	VOUCHER_OPERATE_ADDKEY       = "0"  //添加公钥
	VOUCHER_OPERATE_CREATE       = "1"  //创建
	VOUCHER_OPERATE_DEPLOY       = "2"  //发布
	VOUCHER_OPERATE_START        = "3"  //启动
	VOUCHER_OPERATE_PAUSE        = "4"  //停止
	VOUCHER_OPERATE_HASH_ENABLE  = "5"  //hash同意
	VOUCHER_OPERATE_HASH_DISABLE = "6"  //hash拒绝
	VOUCHER_OPERATE_HASH_LIST    = "7"  //hash list 查询
	VOUCHER_OPERATE_TOKEN_ADD    = "8"  //token 添加
	VOUCHER_OPERATE_TOKEN_DEL    = "9"  //token 删除
	VOUCHER_OPERATE_TOKEN_LIST   = "10" //token list 查询
	VOUCHER_OPERATE_COIN         = "11" //coin 操作
)

type BoxRecord struct {
	Allowed bool
	Address string
	Amount  int64
	Hash    string
}

type SignHashRecord struct {
	SignTmpHash common.Hash
}

type ApplyRecord struct {
	Signflow     string
	SignTmplHash common.Hash
	TxHash       common.Hash
	Amount       *big.Int
	Fee          *big.Int
	To           common.Address
	IsERC20      bool
}

//grpc stream
type GrpcStream struct {
	Type           string
	BlockNumber    uint64 //区块号
	AppId          string //申请人
	Hash           common.Hash
	WdHash         common.Hash
	TxHash         string
	Amount         *big.Int
	Fee            *big.Int
	Account        string
	From           string
	To             string
	Category       *big.Int
	Flow           string //原始内容
	Sign           string //签名信息
	WdFlow         string //提现原始数据
	Status         string
	VoucherOperate *Operate
	ApplyTime      time.Time //申请时间
	TokenList      []*TokenInfo
	SignInfos      []*SignInfo
}

type TokenInfo struct {
	TokenName    string
	Decimals     int64
	ContractAddr string
	Category     int64
}

type SignInfo struct {
	AppId string
	Sign  string
}

type Ecr20Record struct {
	Type     string
	Hash     common.Hash
	WdHash   common.Hash
	Token    common.Address
	To       common.Address
	Amount   *big.Int
	Category *big.Int //转账类型
}

//上公链操作
var Ecr20RecordChan chan *Ecr20Record = make(chan *Ecr20Record, CHAN_MAX_SIZE)

type BtcRecord struct {
	Type string
	//Info 	  	string
	From      string
	WdHash    common.Hash
	Txid      string
	Handlefee *big.Int
	To        string
	Amount    *big.Int
	//Category *big.Int //转账类型
}

var BtcRecordChan chan *BtcRecord = make(chan *BtcRecord, CHAN_MAX_SIZE)

//
func (box *BoxRecord) EncodeToBytes() ([]byte, error) {
	boxJson, err := json.Marshal(box)
	if err != nil {
		return nil, err
	}
	return boxJson, nil
}

func DecodeToBoxRecord(data []byte) (*BoxRecord, error) {
	var box BoxRecord
	if err := json.Unmarshal(data, &box); err != nil {
		return nil, err
	}
	return &box, nil
}

func (meta MetaKey) Bytes() []byte {
	return []byte(meta)
}

func (stage InitStage) Bytes() []byte {
	return []byte(stage)
}

func (stage InitStage) Equals(other []byte) bool {
	return bytes.Compare(stage, other) == 0
}

const (
	ACCOUNT_NO_USED     = "0"
	ACCOUNT_USED        = "1"
	ACCOUNT_BTC         = "BTC"
	ACCOUNT_BTC_NO_USED = "BTC_0_"
	ACCOUNT_BTC_USED    = "BTC_1_"
	ACCOUNT_BTC_COUNT   = "countbtc"
	ACCOUNT_BTC_IMPORT  = "importbtc"

	//ACCOUNT_NONCE   		= "nonce"
)
const (
	BTC_CURSOR = "CURSORBTC"
)

const (
	BTC_DB_REC        = "DBBTCREC_"  //充值
	BTC_DB_WITHDRAW   = "DBBTCWD_"   //WITHDRAW			提现，提币
	BTC_DB_WITHDRAW_0 = "DBBTCWD_0_" //WITHDRAW			提现，提币 交易未被确认
	BTC_DB_WITHDRAW_1 = "DBBTCWD_1_" //WITHDRAW			提现，提币 交易被确认
)
const (
	ACCOUNT_TYPE_ETH = "0"
	ACCOUNT_TYPE_BTC = "1"
)

//审批流状态
const (
	HASH_STATUS_APPLY   = "1" //申请
	HASH_STATUS_ENABLE  = "2" //确认
	HASH_STATUS_DISABLE = "3" //禁用
)

//确认、禁用请求
//const (
//	PASS_STATUS_0 = "0"
//	PASS_STATUS_1 = "1"
//	PASS_STATUS_2 = "2"
//	PASS_STATUS_3 = "3"
//)

//TODO
//category
const (
	COIN_NAME_BTC           = "BTC"
	COIN_NAME_ETH           = "ETH"
	CATEGORY_BTC      int64 = 0
	CATEGORY_ETH      int64 = 1
	COIN_DECIMALS_BTC int64 = 8
	COIN_DECIMALS_ETH int64 = 18
)

var AccountUsedMap map[string]string = make(map[string]string)

var AccountUsedChan chan common.Address = make(chan common.Address, CHAN_MAX_SIZE)

var ReportedChan chan *GrpcStream = make(chan *GrpcStream, CHAN_MAX_SIZE)

//[hash]txid
var UnconfirmedTxidList map[string]string = make(map[string]string)

//BTC
type BTCAddrState struct {
	Addr     string
	IsUsed   bool
	Index    uint32
	IsImport bool
}

var AccountMapBtc map[string]*BTCAddrState = make(map[string]*BTCAddrState)

var AccountImpChanBtc chan *BTCAddrState = make(chan *BTCAddrState, CHAN_MAX_SIZE)

var AccountUsedChanBtc chan string = make(chan string, CHAN_MAX_SIZE)

//私钥-签名机操作
type Operate struct {
	Type         string
	AppId        string //appid
	AppName      string //app别名
	Hash         string
	Password     string
	ReqIpPort    string
	Role         string
	PublicKey    string //加密后公钥
	TokenName    string
	Decimals     int64
	ContractAddr string
	CoinCategory int64  //币种分类
	CoinUsed     bool   //币种使用
	Sign         string //签名
}

var OperateChan chan *Operate = make(chan *Operate, CHAN_MAX_SIZE)

//签名机状态
type VoucherStatus struct {
	ServerStatus    int              //系统状态
	Status          int              //错误码状态
	Total           int              //密钥数量
	HashCount       int              //hash数量
	TokenCount      int              //token数量
	Address         string           //账户地址
	ContractAddress string           //合约地址
	BtcAddress      string           //比特币地址
	D               int              //随机数
	NodesAuthorized []NodeAuthorized //授权情况
	KeyStoreStatus  []KeyStoreStatu  //公约添加状态
	CoinStatus      []CoinStatu      //币种使用状态
}

type NodeAuthorized struct {
	ApplyerId  string
	Authorized bool
}

type KeyStoreStatu struct {
	ApplyerId   string
	ApplyerName string
}

type CoinStatu struct {
	Name     string
	Category int64
	Decimals int64
	Used     bool
}

//实时状态
var RealTimeStatus *VoucherStatus = &VoucherStatus{}

//hash审批流配置模版信息
type HashFlow struct {
	Flow_name     string
	Single_limit  string
	Approval_info []ApprovalInfo
}

type ApprovalInfo struct {
	Require   int64
	Total     int64
	Approvers []Approver
}

type Approver struct {
	Account        string
	ItemType       int64
	Pub_key        string
	App_account_id string
}

type EmployeeFlow struct {
	Appid string
	Sign  string
}
