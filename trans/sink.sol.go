// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package trans

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// OracleABI is the input ABI used to generate the binding from.
const OracleABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalEnabledNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"disableSigner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"isSigner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"indexOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"boss\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"addSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// OracleBin is the compiled bytecode used for deploying new contracts.
const OracleBin = `0x6060604052341561000f57600080fd5b60008054600160a060020a033316600160a060020a03199091161790556105a88061003b6000396000f3006060604052600436106100825763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306661abd81146100875780630921953e146100ac57806327086336146100bf5780637df73e27146100f257806391ac7e6514610111578063c772af391461014b578063eb12d61e1461017a575b600080fd5b341561009257600080fd5b61009a61019b565b60405190815260200160405180910390f35b34156100b757600080fd5b61009a6101a2565b34156100ca57600080fd5b6100de600160a060020a03600435166101f8565b604051901515815260200160405180910390f35b34156100fd57600080fd5b6100de600160a060020a0360043516610291565b341561011c57600080fd5b6101276004356102f8565b604051600160a060020a039092168252151560208201526040908101905180910390f35b341561015657600080fd5b61015e610357565b604051600160a060020a03909116815260200160405180910390f35b341561018557600080fd5b610199600160a060020a0360043516610366565b005b6002545b90565b60025460009060015b6002548110156101ee5760028054829081106101c357fe5b60009182526020909120015460a060020a900460ff1615156101e6576001820391505b6001016101ab565b5060001901919050565b60008054819033600160a060020a0390811691161461021657600080fd5b50600160a060020a038216600090815260016020526040902054801561028657600060028281548110151561024757fe5b6000918252602090912001805491151560a060020a0274ff0000000000000000000000000000000000000000199092169190911790556001915061028b565b600091505b50919050565b600160a060020a038116600090815260016020526040812054158015906102f25750600160a060020a0382166000908152600160205260409020546002805490919081106102db57fe5b60009182526020909120015460a060020a900460ff165b92915050565b60008060008060028581548110151561030d57fe5b60009182526020909120015460028054600160a060020a039092169350908690811061033557fe5b6000918252602090912001549193505060ff60a060020a909104169050915091565b600054600160a060020a031681565b60008054819033600160a060020a0390811691161461038457600080fd5b600254151561045057600161039a60028261051d565b50600080805260016020527fa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb4955604080519081016040526000808252602082018190526002805490919081106103ec57fe5b60009182526020909120018151815473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03919091161781556020820151815490151560a060020a0274ff000000000000000000000000000000000000000019909116179055505b600160a060020a03831660009081526001602052604090205491508115156104a65760028054600160a060020a038516600090815260016020819052604090912082905590935083916104a491830161051d565b505b60028054839081106104b457fe5b60009182526020808320909101805474ff000000000000000000000000000000000000000019600160a060020a0390971673ffffffffffffffffffffffffffffffffffffffff1990911681179690961660a060020a179055938152600190935250604090912055565b81548183558181151161054157600083815260209020610541918101908301610546565b505050565b61019f91905b8082111561057857805474ffffffffffffffffffffffffffffffffffffffffff1916815560010161054c565b50905600a165627a7a72305820e3889acff26dcee79862518ad1330a2c070a106e6b21cd1065bef4857dc382e80029`

// DeployOracle deploys a new Ethereum contract, binding an instance of Oracle to it.
func DeployOracle(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Oracle, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OracleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// Boss is a free data retrieval call binding the contract method 0xc772af39.
//
// Solidity: function boss() constant returns(address)
func (_Oracle *OracleCaller) Boss(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "boss")
	return *ret0, err
}

// Boss is a free data retrieval call binding the contract method 0xc772af39.
//
// Solidity: function boss() constant returns(address)
func (_Oracle *OracleSession) Boss() (common.Address, error) {
	return _Oracle.Contract.Boss(&_Oracle.CallOpts)
}

// Boss is a free data retrieval call binding the contract method 0xc772af39.
//
// Solidity: function boss() constant returns(address)
func (_Oracle *OracleCallerSession) Boss() (common.Address, error) {
	return _Oracle.Contract.Boss(&_Oracle.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() constant returns(uint256)
func (_Oracle *OracleCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "count")
	return *ret0, err
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() constant returns(uint256)
func (_Oracle *OracleSession) Count() (*big.Int, error) {
	return _Oracle.Contract.Count(&_Oracle.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() constant returns(uint256)
func (_Oracle *OracleCallerSession) Count() (*big.Int, error) {
	return _Oracle.Contract.Count(&_Oracle.CallOpts)
}

// IndexOf is a free data retrieval call binding the contract method 0x91ac7e65.
//
// Solidity: function indexOf(idx uint256) constant returns(address, bool)
func (_Oracle *OracleCaller) IndexOf(opts *bind.CallOpts, idx *big.Int) (common.Address, bool, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(bool)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Oracle.contract.Call(opts, out, "indexOf", idx)
	return *ret0, *ret1, err
}

// IndexOf is a free data retrieval call binding the contract method 0x91ac7e65.
//
// Solidity: function indexOf(idx uint256) constant returns(address, bool)
func (_Oracle *OracleSession) IndexOf(idx *big.Int) (common.Address, bool, error) {
	return _Oracle.Contract.IndexOf(&_Oracle.CallOpts, idx)
}

// IndexOf is a free data retrieval call binding the contract method 0x91ac7e65.
//
// Solidity: function indexOf(idx uint256) constant returns(address, bool)
func (_Oracle *OracleCallerSession) IndexOf(idx *big.Int) (common.Address, bool, error) {
	return _Oracle.Contract.IndexOf(&_Oracle.CallOpts, idx)
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(signer address) constant returns(bool)
func (_Oracle *OracleCaller) IsSigner(opts *bind.CallOpts, signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "isSigner", signer)
	return *ret0, err
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(signer address) constant returns(bool)
func (_Oracle *OracleSession) IsSigner(signer common.Address) (bool, error) {
	return _Oracle.Contract.IsSigner(&_Oracle.CallOpts, signer)
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(signer address) constant returns(bool)
func (_Oracle *OracleCallerSession) IsSigner(signer common.Address) (bool, error) {
	return _Oracle.Contract.IsSigner(&_Oracle.CallOpts, signer)
}

// TotalEnabledNodes is a free data retrieval call binding the contract method 0x0921953e.
//
// Solidity: function totalEnabledNodes() constant returns(uint256)
func (_Oracle *OracleCaller) TotalEnabledNodes(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "totalEnabledNodes")
	return *ret0, err
}

// TotalEnabledNodes is a free data retrieval call binding the contract method 0x0921953e.
//
// Solidity: function totalEnabledNodes() constant returns(uint256)
func (_Oracle *OracleSession) TotalEnabledNodes() (*big.Int, error) {
	return _Oracle.Contract.TotalEnabledNodes(&_Oracle.CallOpts)
}

// TotalEnabledNodes is a free data retrieval call binding the contract method 0x0921953e.
//
// Solidity: function totalEnabledNodes() constant returns(uint256)
func (_Oracle *OracleCallerSession) TotalEnabledNodes() (*big.Int, error) {
	return _Oracle.Contract.TotalEnabledNodes(&_Oracle.CallOpts)
}

// AddSigner is a paid mutator transaction binding the contract method 0xeb12d61e.
//
// Solidity: function addSigner(signer address) returns()
func (_Oracle *OracleTransactor) AddSigner(opts *bind.TransactOpts, signer common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "addSigner", signer)
}

// AddSigner is a paid mutator transaction binding the contract method 0xeb12d61e.
//
// Solidity: function addSigner(signer address) returns()
func (_Oracle *OracleSession) AddSigner(signer common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.AddSigner(&_Oracle.TransactOpts, signer)
}

// AddSigner is a paid mutator transaction binding the contract method 0xeb12d61e.
//
// Solidity: function addSigner(signer address) returns()
func (_Oracle *OracleTransactorSession) AddSigner(signer common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.AddSigner(&_Oracle.TransactOpts, signer)
}

// DisableSigner is a paid mutator transaction binding the contract method 0x27086336.
//
// Solidity: function disableSigner(signer address) returns(bool)
func (_Oracle *OracleTransactor) DisableSigner(opts *bind.TransactOpts, signer common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "disableSigner", signer)
}

// DisableSigner is a paid mutator transaction binding the contract method 0x27086336.
//
// Solidity: function disableSigner(signer address) returns(bool)
func (_Oracle *OracleSession) DisableSigner(signer common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.DisableSigner(&_Oracle.TransactOpts, signer)
}

// DisableSigner is a paid mutator transaction binding the contract method 0x27086336.
//
// Solidity: function disableSigner(signer address) returns(bool)
func (_Oracle *OracleTransactorSession) DisableSigner(signer common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.DisableSigner(&_Oracle.TransactOpts, signer)
}

// SinkABI is the input ABI used to generate the binding from.
const SinkABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"enable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"addHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOracle\",\"type\":\"address\"}],\"name\":\"changeOracle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"available\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"txHash\",\"type\":\"bytes32\"}],\"name\":\"txExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\"},{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"category\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"disable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"ref\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"lastConfirmed\",\"type\":\"address\"}],\"name\":\"SignflowAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"lastConfirmed\",\"type\":\"address\"}],\"name\":\"SignflowEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"lastConfirmed\",\"type\":\"address\"}],\"name\":\"SignflowDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"category\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"lastConfirmed\",\"type\":\"address\"}],\"name\":\"WithdrawApplied\",\"type\":\"event\"}]"

// SinkBin is the compiled bytecode used for deploying new contracts.
const SinkBin = `0x606060405234156200001057600080fd5b604051602080620019a08339810160405280805160008054600160a060020a031916600160a060020a038316178155909250620000629150808080808080640100000000620000778102620005281704565b50506003805460ff1916600117905562000ad3565b600086815260016020526040812054818080831515620000fb5760008c6003811115620000a057fe5b14620000b05760009450620004b6565b8a158015620000c1575060035460ff165b15620000d15760009450620004b6565b6002805460008d81526001602081905260409091208290559091620000f891830162000a29565b93505b60028054859081106200010a57fe5b60009182526020822060049091020193508c60038111156200012857fe5b141562000251576001830154610100900460ff1615156200020d578a835560018301805461ff00191661010017905562000174838d8c8c8a640100000000620008f5620004c582021704565b1515620001855760009450620004b6565b6200019d64010000000062000af3620006e282021704565b620001b7848e64010000000062000b9e620007ad82021704565b106200020357600080516020620019808339815191528b33604051918252600160a060020a031660208201526040908101905180910390a160018301805462ff00001916620100001790555b60019450620004b6565b600183015462010000900460ff16806200024157506200023f838d8c8c8a640100000000620008f5620004c582021704565b155b15620002515760009450620004b6565b6200026e838d8c8c8a64010000000062000bd2620007e282021704565b9150915081156200028257809450620004b6565b60008c60038111156200029157fe5b1415620002e357600080516020620019808339815191528b33604051918252600160a060020a031660208201526040908101905180910390a160018301805462ff0000191662010000179055620004b1565b60018c6003811115620002f257fe5b14156200036c576001838101805460ff1916909117905562000323838d64010000000062000cc96200091882021704565b7f50177799234754a6d6af99e5ab43b5679c202f4058d342099bfb35acdfa1a8678b33604051918252600160a060020a031660208201526040908101905180910390a1620004b1565b60028c60038111156200037b57fe5b1415620003f15760018301805460ff19169055620003a8838d64010000000062000cc96200091882021704565b7f484f49c7a40838d935f9cd616461fad6033bb6f7fa4491fbc72941d77671f09f8b33604051918252600160a060020a031660208201526040908101905180910390a1620004b1565b60038c60038111156200040057fe5b1415620004a75760008a81526003840160205260409020805460ff191660011790556200043c838d64010000000062000cc96200091882021704565b898b7f7f508fd15756f38a4426383f9ef243dfccdfdff0a528e755b113ef8bef2c5c2e8b8b8b8b336040519485526020850193909352600160a060020a0391821660408086019190915260608501919091529116608083015260a0909101905180910390a3620004b1565b60009450620004b6565b600194505b50505050979650505050505050565b6000806000876002016000886003811115620004dd57fe5b60ff1681526020810191909152604001600020915060018760038111156200050157fe5b148015620005135750600188015460ff165b806200053c575060028760038111156200052957fe5b1480156200053c5750600188015460ff16155b8062000565575060038760038111156200055257fe5b148015620005655750600188015460ff16155b15620005755760009250620006d7565b6200058f888864010000000062000d89620009de82021704565b156200059f5760009250620006d7565b6003876003811115620005ae57fe5b148015620005be57506002820154155b15620005df5760038201859055600582018690556004820184905562000639565b6003876003811115620005ee57fe5b14801562000600575060008260020154115b1562000639578482600301541480156200061d5750600582015486145b80156200062d5750838260040154145b15156200063957600080fd5b600160a060020a0333166000908152602083905260408120805460ff1916600190811790915560028401805490910190558760038111156200067757fe5b14620006d2576001808301805491620006939190830162000a5d565b9050338260010182815481101515620006a857fe5b60009182526020909120018054600160a060020a031916600160a060020a03929092169190911790555b600192505b505095945050505050565b6000805481908190600160a060020a0316630921953e82604051602001526040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b15156200074857600080fd5b6102c65a03f115156200075a57600080fd5b50505060405180519250506001821480620007755750816002145b156200078457819250620007a8565b506002810680156200079e578060028304019250620007a8565b6002820460010192505b505090565b600080836002016000846003811115620007c357fe5b60ff168152602081019190915260400160002060020154949350505050565b6000806003866003811115620007f457fe5b148015620008125750600085815260038801602052604090205460ff165b156200082557506001905060006200090e565b60008660038111156200083457fe5b141580156200086357506001870154610100900460ff161580620008635750600187015462010000900460ff16155b156200087657506001905060006200090e565b60008660038111156200088557fe5b14158015620008ae5750620008ac8787878787640100000000620008f5620004c582021704565b155b15620008c157506001905060006200090e565b620008d964010000000062000af3620006e282021704565b620008f3888864010000000062000b9e620007ad82021704565b101562000906575060019050806200090e565b506000905060015b9550959350505050565b60008060008460020160008560038111156200093057fe5b60ff1660ff168152602001908152602001600020925082600101805490509050600091505b80821015620009ad576001830180548491600091859081106200097457fe5b6000918252602080832090910154600160a060020a031683528201929092526040019020805460ff191690556001919091019062000955565b600060028401819055620009c5600185018262000a5d565b5050600060038301819055600590920191909155505050565b600080836002016000846003811115620009f457fe5b60ff90811682526020808301939093526040918201600090812033600160a060020a0316825290935291205416949350505050565b81548183558181151162000a585760040281600402836000526020600020918201910162000a58919062000a84565b505050565b81548183558181151162000a585760008381526020902062000a5891810190830162000ab6565b62000ab391905b8082111562000aaf576000815560018101805462ffffff1916905560040162000a8b565b5090565b90565b62000ab391905b8082111562000aaf576000815560010162000abd565b610e9d8062000ae36000396000f3006060604052600436106100695763ffffffff60e060020a6000350416631059171e811461006e57806343e08ad11461009857806347c421b5146100ae5780636932854f146100cf578063a8022bfa146100ff578063ca5031ab14610118578063cf751d8114610146575b600080fd5b341561007957600080fd5b61008460043561015c565b604051901515815260200160405180910390f35b34156100a357600080fd5b6100846004356101f2565b34156100b957600080fd5b6100cd600160a060020a0360043516610281565b005b34156100da57600080fd5b6100e560043561032e565b604051918252151560208201526040908101905180910390f35b341561010a57600080fd5b6100846004356024356103a4565b341561012357600080fd5b610084600435602435604435600160a060020a036064351660843560a435610402565b341561015157600080fd5b61008460043561049c565b60008054600160a060020a0316637df73e2733836040516020015260405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b15156101b557600080fd5b6102c65a03f115156101c657600080fd5b5050506040518051905015156101db57600080fd5b6101ec600183600080808080610528565b92915050565b60008054600160a060020a0316637df73e2733836040516020015260405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b151561024b57600080fd5b6102c65a03f1151561025c57600080fd5b50505060405180519050151561027157600080fd5b6101ec6000838180808080610528565b60008054600160a060020a033381169291169063c772af3990604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b15156102cf57600080fd5b6102c65a03f115156102e057600080fd5b50505060405180519050600160a060020a03161415156102ff57600080fd5b6000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b6000818152600160205260408120548190801515610352576000925082915061039e565b600280548290811061036057fe5b90600052602060002090600402016000015460028281548110151561038157fe5b600091825260209091206001600490920201015490935060ff1691505b50915091565b600082815260016020526040812054818115156103c457600092506103fa565b60028054839081106103d257fe5b600091825260208083208784526003600490930201918201905260409091205460ff16935090505b505092915050565b60008054600160a060020a0316637df73e2733836040516020015260405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b151561045b57600080fd5b6102c65a03f1151561046c57600080fd5b50505060405180519050151561048157600080fd5b6104916003848989898988610528565b979650505050505050565b60008054600160a060020a0316637df73e2733836040516020015260405160e060020a63ffffffff8416028152600160a060020a039091166004820152602401602060405180830381600087803b15156104f557600080fd5b6102c65a03f1151561050657600080fd5b50505060405180519050151561051b57600080fd5b6101ec6002836000808080805b6000868152600160205260408120548180808315156105a35760008c600381111561054f57fe5b1461055d57600094506108e6565b8a15801561056d575060035460ff165b1561057b57600094506108e6565b6002805460008d815260016020819052604090912082905590916105a0918301610dd3565b93505b60028054859081106105b157fe5b60009182526020822060049091020193508c60038111156105ce57fe5b14156106bf576001830154610100900460ff16151561068e578a835560018301805461ff001916610100179055610608838d8c8c8a6108f5565b151561061757600094506108e6565b61061f610af3565b610629848e610b9e565b10610685577f9f3f4c1672a4880364b07219cd9428dbc8a88774f53b12b98fae406c5a30ee5c8b33604051918252600160a060020a031660208201526040908101905180910390a160018301805462ff00001916620100001790555b600194506108e6565b600183015462010000900460ff16806106b157506106af838d8c8c8a6108f5565b155b156106bf57600094506108e6565b6106cc838d8c8c8a610bd2565b9150915081156106de578094506108e6565b60008c60038111156106ec57fe5b141561074d577f9f3f4c1672a4880364b07219cd9428dbc8a88774f53b12b98fae406c5a30ee5c8b33604051918252600160a060020a031660208201526040908101905180910390a160018301805462ff00001916620100001790556108e1565b60018c600381111561075b57fe5b14156107c3576001838101805460ff1916909117905561077b838d610cc9565b7f50177799234754a6d6af99e5ab43b5679c202f4058d342099bfb35acdfa1a8678b33604051918252600160a060020a031660208201526040908101905180910390a16108e1565b60028c60038111156107d157fe5b14156108355760018301805460ff191690556107ed838d610cc9565b7f484f49c7a40838d935f9cd616461fad6033bb6f7fa4491fbc72941d77671f09f8b33604051918252600160a060020a031660208201526040908101905180910390a16108e1565b60038c600381111561084357fe5b14156108d85760008a81526003840160205260409020805460ff1916600117905561086e838d610cc9565b898b7f7f508fd15756f38a4426383f9ef243dfccdfdff0a528e755b113ef8bef2c5c2e8b8b8b8b336040519485526020850193909352600160a060020a0391821660408086019190915260608501919091529116608083015260a0909101905180910390a36108e1565b600094506108e6565b600194505b50505050979650505050505050565b600080600087600201600088600381111561090c57fe5b60ff16815260208101919091526040016000209150600187600381111561092f57fe5b1480156109405750600188015460ff165b806109665750600287600381111561095457fe5b1480156109665750600188015460ff16155b8061098c5750600387600381111561097a57fe5b14801561098c5750600188015460ff16155b1561099a5760009250610ae8565b6109a48888610d89565b156109b25760009250610ae8565b60038760038111156109c057fe5b1480156109cf57506002820154155b156109ee57600382018590556005820186905560048201849055610a42565b60038760038111156109fc57fe5b148015610a0d575060008260020154115b15610a4257848260030154148015610a285750600582015486145b8015610a375750838260040154145b1515610a4257600080fd5b600160a060020a0333166000908152602083905260408120805460ff191660019081179091556002840180549091019055876003811115610a7f57fe5b14610ae3576001808301805491610a9891908301610e04565b9050338260010182815481101515610aac57fe5b6000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555b600192505b505095945050505050565b6000805481908190600160a060020a0316630921953e82604051602001526040518163ffffffff1660e060020a028152600401602060405180830381600087803b1515610b3f57600080fd5b6102c65a03f11515610b5057600080fd5b50505060405180519250506001821480610b6a5750816002145b15610b7757819250610b99565b50600281068015610b8f578060028304019250610b99565b6002820460010192505b505090565b600080836002016000846003811115610bb357fe5b60ff168152602081019190915260400160002060020154949350505050565b6000806003866003811115610be357fe5b148015610c005750600085815260038801602052604090205460ff165b15610c115750600190506000610cbf565b6000866003811115610c1f57fe5b14158015610c4c57506001870154610100900460ff161580610c4c5750600187015462010000900460ff16155b15610c5d5750600190506000610cbf565b6000866003811115610c6b57fe5b14158015610c835750610c8187878787876108f5565b155b15610c945750600190506000610cbf565b610c9c610af3565b610ca68888610b9e565b1015610cb757506001905080610cbf565b506000905060015b9550959350505050565b6000806000846002016000856003811115610ce057fe5b60ff1660ff168152602001908152602001600020925082600101805490509050600091505b80821015610d5a57600183018054849160009185908110610d2257fe5b6000918252602080832090910154600160a060020a031683528201929092526040019020805460ff1916905560019190910190610d05565b600060028401819055610d706001850182610e04565b5050600060038301819055600590920191909155505050565b600080836002016000846003811115610d9e57fe5b60ff90811682526020808301939093526040918201600090812033600160a060020a0316825290935291205416949350505050565b815481835581811511610dff57600402816004028360005260206000209182019101610dff9190610e28565b505050565b815481835581811511610dff57600083815260209020610dff918101908301610e57565b610e5491905b80821115610e50576000815560018101805462ffffff19169055600401610e2e565b5090565b90565b610e5491905b80821115610e505760008155600101610e5d5600a165627a7a723058205fe60cbcbdd3f23e67517d6f64f7cef33821eeeb90029356737d33c80abdb0db00299f3f4c1672a4880364b07219cd9428dbc8a88774f53b12b98fae406c5a30ee5c`

// DeploySink deploys a new Ethereum contract, binding an instance of Sink to it.
func DeploySink(auth *bind.TransactOpts, backend bind.ContractBackend, ref common.Address) (common.Address, *types.Transaction, *Sink, error) {
	parsed, err := abi.JSON(strings.NewReader(SinkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SinkBin), backend, ref)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Sink{SinkCaller: SinkCaller{contract: contract}, SinkTransactor: SinkTransactor{contract: contract}, SinkFilterer: SinkFilterer{contract: contract}}, nil
}

// Sink is an auto generated Go binding around an Ethereum contract.
type Sink struct {
	SinkCaller     // Read-only binding to the contract
	SinkTransactor // Write-only binding to the contract
	SinkFilterer   // Log filterer for contract events
}

// SinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type SinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SinkSession struct {
	Contract     *Sink             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SinkCallerSession struct {
	Contract *SinkCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SinkTransactorSession struct {
	Contract     *SinkTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type SinkRaw struct {
	Contract *Sink // Generic contract binding to access the raw methods on
}

// SinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SinkCallerRaw struct {
	Contract *SinkCaller // Generic read-only contract binding to access the raw methods on
}

// SinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SinkTransactorRaw struct {
	Contract *SinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSink creates a new instance of Sink, bound to a specific deployed contract.
func NewSink(address common.Address, backend bind.ContractBackend) (*Sink, error) {
	contract, err := bindSink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Sink{SinkCaller: SinkCaller{contract: contract}, SinkTransactor: SinkTransactor{contract: contract}, SinkFilterer: SinkFilterer{contract: contract}}, nil
}

// NewSinkCaller creates a new read-only instance of Sink, bound to a specific deployed contract.
func NewSinkCaller(address common.Address, caller bind.ContractCaller) (*SinkCaller, error) {
	contract, err := bindSink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SinkCaller{contract: contract}, nil
}

// NewSinkTransactor creates a new write-only instance of Sink, bound to a specific deployed contract.
func NewSinkTransactor(address common.Address, transactor bind.ContractTransactor) (*SinkTransactor, error) {
	contract, err := bindSink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SinkTransactor{contract: contract}, nil
}

// NewSinkFilterer creates a new log filterer instance of Sink, bound to a specific deployed contract.
func NewSinkFilterer(address common.Address, filterer bind.ContractFilterer) (*SinkFilterer, error) {
	contract, err := bindSink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SinkFilterer{contract: contract}, nil
}

// bindSink binds a generic wrapper to an already deployed contract.
func bindSink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sink *SinkRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Sink.Contract.SinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sink *SinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sink.Contract.SinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sink *SinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sink.Contract.SinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sink *SinkCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Sink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sink *SinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sink *SinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sink.Contract.contract.Transact(opts, method, params...)
}

// Available is a free data retrieval call binding the contract method 0x6932854f.
//
// Solidity: function available(hash bytes32) constant returns(bytes32, bool)
func (_Sink *SinkCaller) Available(opts *bind.CallOpts, hash [32]byte) ([32]byte, bool, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Sink.contract.Call(opts, out, "available", hash)
	return *ret0, *ret1, err
}

// Available is a free data retrieval call binding the contract method 0x6932854f.
//
// Solidity: function available(hash bytes32) constant returns(bytes32, bool)
func (_Sink *SinkSession) Available(hash [32]byte) ([32]byte, bool, error) {
	return _Sink.Contract.Available(&_Sink.CallOpts, hash)
}

// Available is a free data retrieval call binding the contract method 0x6932854f.
//
// Solidity: function available(hash bytes32) constant returns(bytes32, bool)
func (_Sink *SinkCallerSession) Available(hash [32]byte) ([32]byte, bool, error) {
	return _Sink.Contract.Available(&_Sink.CallOpts, hash)
}

// TxExists is a free data retrieval call binding the contract method 0xa8022bfa.
//
// Solidity: function txExists(hash bytes32, txHash bytes32) constant returns(bool)
func (_Sink *SinkCaller) TxExists(opts *bind.CallOpts, hash [32]byte, txHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Sink.contract.Call(opts, out, "txExists", hash, txHash)
	return *ret0, err
}

// TxExists is a free data retrieval call binding the contract method 0xa8022bfa.
//
// Solidity: function txExists(hash bytes32, txHash bytes32) constant returns(bool)
func (_Sink *SinkSession) TxExists(hash [32]byte, txHash [32]byte) (bool, error) {
	return _Sink.Contract.TxExists(&_Sink.CallOpts, hash, txHash)
}

// TxExists is a free data retrieval call binding the contract method 0xa8022bfa.
//
// Solidity: function txExists(hash bytes32, txHash bytes32) constant returns(bool)
func (_Sink *SinkCallerSession) TxExists(hash [32]byte, txHash [32]byte) (bool, error) {
	return _Sink.Contract.TxExists(&_Sink.CallOpts, hash, txHash)
}

// AddHash is a paid mutator transaction binding the contract method 0x43e08ad1.
//
// Solidity: function addHash(hash bytes32) returns(bool)
func (_Sink *SinkTransactor) AddHash(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _Sink.contract.Transact(opts, "addHash", hash)
}

// AddHash is a paid mutator transaction binding the contract method 0x43e08ad1.
//
// Solidity: function addHash(hash bytes32) returns(bool)
func (_Sink *SinkSession) AddHash(hash [32]byte) (*types.Transaction, error) {
	return _Sink.Contract.AddHash(&_Sink.TransactOpts, hash)
}

// AddHash is a paid mutator transaction binding the contract method 0x43e08ad1.
//
// Solidity: function addHash(hash bytes32) returns(bool)
func (_Sink *SinkTransactorSession) AddHash(hash [32]byte) (*types.Transaction, error) {
	return _Sink.Contract.AddHash(&_Sink.TransactOpts, hash)
}

// Approve is a paid mutator transaction binding the contract method 0xca5031ab.
//
// Solidity: function approve(txHash bytes32, amount uint256, fee uint256, recipient address, hash bytes32, category uint256) returns(bool)
func (_Sink *SinkTransactor) Approve(opts *bind.TransactOpts, txHash [32]byte, amount *big.Int, fee *big.Int, recipient common.Address, hash [32]byte, category *big.Int) (*types.Transaction, error) {
	return _Sink.contract.Transact(opts, "approve", txHash, amount, fee, recipient, hash, category)
}

// Approve is a paid mutator transaction binding the contract method 0xca5031ab.
//
// Solidity: function approve(txHash bytes32, amount uint256, fee uint256, recipient address, hash bytes32, category uint256) returns(bool)
func (_Sink *SinkSession) Approve(txHash [32]byte, amount *big.Int, fee *big.Int, recipient common.Address, hash [32]byte, category *big.Int) (*types.Transaction, error) {
	return _Sink.Contract.Approve(&_Sink.TransactOpts, txHash, amount, fee, recipient, hash, category)
}

// Approve is a paid mutator transaction binding the contract method 0xca5031ab.
//
// Solidity: function approve(txHash bytes32, amount uint256, fee uint256, recipient address, hash bytes32, category uint256) returns(bool)
func (_Sink *SinkTransactorSession) Approve(txHash [32]byte, amount *big.Int, fee *big.Int, recipient common.Address, hash [32]byte, category *big.Int) (*types.Transaction, error) {
	return _Sink.Contract.Approve(&_Sink.TransactOpts, txHash, amount, fee, recipient, hash, category)
}

// ChangeOracle is a paid mutator transaction binding the contract method 0x47c421b5.
//
// Solidity: function changeOracle(newOracle address) returns()
func (_Sink *SinkTransactor) ChangeOracle(opts *bind.TransactOpts, newOracle common.Address) (*types.Transaction, error) {
	return _Sink.contract.Transact(opts, "changeOracle", newOracle)
}

// ChangeOracle is a paid mutator transaction binding the contract method 0x47c421b5.
//
// Solidity: function changeOracle(newOracle address) returns()
func (_Sink *SinkSession) ChangeOracle(newOracle common.Address) (*types.Transaction, error) {
	return _Sink.Contract.ChangeOracle(&_Sink.TransactOpts, newOracle)
}

// ChangeOracle is a paid mutator transaction binding the contract method 0x47c421b5.
//
// Solidity: function changeOracle(newOracle address) returns()
func (_Sink *SinkTransactorSession) ChangeOracle(newOracle common.Address) (*types.Transaction, error) {
	return _Sink.Contract.ChangeOracle(&_Sink.TransactOpts, newOracle)
}

// Disable is a paid mutator transaction binding the contract method 0xcf751d81.
//
// Solidity: function disable(hash bytes32) returns(bool)
func (_Sink *SinkTransactor) Disable(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _Sink.contract.Transact(opts, "disable", hash)
}

// Disable is a paid mutator transaction binding the contract method 0xcf751d81.
//
// Solidity: function disable(hash bytes32) returns(bool)
func (_Sink *SinkSession) Disable(hash [32]byte) (*types.Transaction, error) {
	return _Sink.Contract.Disable(&_Sink.TransactOpts, hash)
}

// Disable is a paid mutator transaction binding the contract method 0xcf751d81.
//
// Solidity: function disable(hash bytes32) returns(bool)
func (_Sink *SinkTransactorSession) Disable(hash [32]byte) (*types.Transaction, error) {
	return _Sink.Contract.Disable(&_Sink.TransactOpts, hash)
}

// Enable is a paid mutator transaction binding the contract method 0x1059171e.
//
// Solidity: function enable(hash bytes32) returns(bool)
func (_Sink *SinkTransactor) Enable(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _Sink.contract.Transact(opts, "enable", hash)
}

// Enable is a paid mutator transaction binding the contract method 0x1059171e.
//
// Solidity: function enable(hash bytes32) returns(bool)
func (_Sink *SinkSession) Enable(hash [32]byte) (*types.Transaction, error) {
	return _Sink.Contract.Enable(&_Sink.TransactOpts, hash)
}

// Enable is a paid mutator transaction binding the contract method 0x1059171e.
//
// Solidity: function enable(hash bytes32) returns(bool)
func (_Sink *SinkTransactorSession) Enable(hash [32]byte) (*types.Transaction, error) {
	return _Sink.Contract.Enable(&_Sink.TransactOpts, hash)
}

// SinkSignflowAddedIterator is returned from FilterSignflowAdded and is used to iterate over the raw logs and unpacked data for SignflowAdded events raised by the Sink contract.
type SinkSignflowAddedIterator struct {
	Event *SinkSignflowAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SinkSignflowAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SinkSignflowAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SinkSignflowAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SinkSignflowAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SinkSignflowAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SinkSignflowAdded represents a SignflowAdded event raised by the Sink contract.
type SinkSignflowAdded struct {
	Hash          [32]byte
	LastConfirmed common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSignflowAdded is a free log retrieval operation binding the contract event 0x9f3f4c1672a4880364b07219cd9428dbc8a88774f53b12b98fae406c5a30ee5c.
//
// Solidity: event SignflowAdded(hash bytes32, lastConfirmed address)
func (_Sink *SinkFilterer) FilterSignflowAdded(opts *bind.FilterOpts) (*SinkSignflowAddedIterator, error) {

	logs, sub, err := _Sink.contract.FilterLogs(opts, "SignflowAdded")
	if err != nil {
		return nil, err
	}
	return &SinkSignflowAddedIterator{contract: _Sink.contract, event: "SignflowAdded", logs: logs, sub: sub}, nil
}

// WatchSignflowAdded is a free log subscription operation binding the contract event 0x9f3f4c1672a4880364b07219cd9428dbc8a88774f53b12b98fae406c5a30ee5c.
//
// Solidity: event SignflowAdded(hash bytes32, lastConfirmed address)
func (_Sink *SinkFilterer) WatchSignflowAdded(opts *bind.WatchOpts, sink chan<- *SinkSignflowAdded) (event.Subscription, error) {

	logs, sub, err := _Sink.contract.WatchLogs(opts, "SignflowAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SinkSignflowAdded)
				if err := _Sink.contract.UnpackLog(event, "SignflowAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SinkSignflowDisabledIterator is returned from FilterSignflowDisabled and is used to iterate over the raw logs and unpacked data for SignflowDisabled events raised by the Sink contract.
type SinkSignflowDisabledIterator struct {
	Event *SinkSignflowDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SinkSignflowDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SinkSignflowDisabled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SinkSignflowDisabled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SinkSignflowDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SinkSignflowDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SinkSignflowDisabled represents a SignflowDisabled event raised by the Sink contract.
type SinkSignflowDisabled struct {
	Hash          [32]byte
	LastConfirmed common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSignflowDisabled is a free log retrieval operation binding the contract event 0x484f49c7a40838d935f9cd616461fad6033bb6f7fa4491fbc72941d77671f09f.
//
// Solidity: event SignflowDisabled(hash bytes32, lastConfirmed address)
func (_Sink *SinkFilterer) FilterSignflowDisabled(opts *bind.FilterOpts) (*SinkSignflowDisabledIterator, error) {

	logs, sub, err := _Sink.contract.FilterLogs(opts, "SignflowDisabled")
	if err != nil {
		return nil, err
	}
	return &SinkSignflowDisabledIterator{contract: _Sink.contract, event: "SignflowDisabled", logs: logs, sub: sub}, nil
}

// WatchSignflowDisabled is a free log subscription operation binding the contract event 0x484f49c7a40838d935f9cd616461fad6033bb6f7fa4491fbc72941d77671f09f.
//
// Solidity: event SignflowDisabled(hash bytes32, lastConfirmed address)
func (_Sink *SinkFilterer) WatchSignflowDisabled(opts *bind.WatchOpts, sink chan<- *SinkSignflowDisabled) (event.Subscription, error) {

	logs, sub, err := _Sink.contract.WatchLogs(opts, "SignflowDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SinkSignflowDisabled)
				if err := _Sink.contract.UnpackLog(event, "SignflowDisabled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SinkSignflowEnabledIterator is returned from FilterSignflowEnabled and is used to iterate over the raw logs and unpacked data for SignflowEnabled events raised by the Sink contract.
type SinkSignflowEnabledIterator struct {
	Event *SinkSignflowEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SinkSignflowEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SinkSignflowEnabled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SinkSignflowEnabled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SinkSignflowEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SinkSignflowEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SinkSignflowEnabled represents a SignflowEnabled event raised by the Sink contract.
type SinkSignflowEnabled struct {
	Hash          [32]byte
	LastConfirmed common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSignflowEnabled is a free log retrieval operation binding the contract event 0x50177799234754a6d6af99e5ab43b5679c202f4058d342099bfb35acdfa1a867.
//
// Solidity: event SignflowEnabled(hash bytes32, lastConfirmed address)
func (_Sink *SinkFilterer) FilterSignflowEnabled(opts *bind.FilterOpts) (*SinkSignflowEnabledIterator, error) {

	logs, sub, err := _Sink.contract.FilterLogs(opts, "SignflowEnabled")
	if err != nil {
		return nil, err
	}
	return &SinkSignflowEnabledIterator{contract: _Sink.contract, event: "SignflowEnabled", logs: logs, sub: sub}, nil
}

// WatchSignflowEnabled is a free log subscription operation binding the contract event 0x50177799234754a6d6af99e5ab43b5679c202f4058d342099bfb35acdfa1a867.
//
// Solidity: event SignflowEnabled(hash bytes32, lastConfirmed address)
func (_Sink *SinkFilterer) WatchSignflowEnabled(opts *bind.WatchOpts, sink chan<- *SinkSignflowEnabled) (event.Subscription, error) {

	logs, sub, err := _Sink.contract.WatchLogs(opts, "SignflowEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SinkSignflowEnabled)
				if err := _Sink.contract.UnpackLog(event, "SignflowEnabled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SinkWithdrawAppliedIterator is returned from FilterWithdrawApplied and is used to iterate over the raw logs and unpacked data for WithdrawApplied events raised by the Sink contract.
type SinkWithdrawAppliedIterator struct {
	Event *SinkWithdrawApplied // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SinkWithdrawAppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SinkWithdrawApplied)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SinkWithdrawApplied)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SinkWithdrawAppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SinkWithdrawAppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SinkWithdrawApplied represents a WithdrawApplied event raised by the Sink contract.
type SinkWithdrawApplied struct {
	Hash          [32]byte
	TxHash        [32]byte
	Amount        *big.Int
	Fee           *big.Int
	Recipient     common.Address
	Category      *big.Int
	LastConfirmed common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterWithdrawApplied is a free log retrieval operation binding the contract event 0x7f508fd15756f38a4426383f9ef243dfccdfdff0a528e755b113ef8bef2c5c2e.
//
// Solidity: event WithdrawApplied(hash indexed bytes32, txHash indexed bytes32, amount uint256, fee uint256, recipient address, category uint256, lastConfirmed address)
func (_Sink *SinkFilterer) FilterWithdrawApplied(opts *bind.FilterOpts, hash [][32]byte, txHash [][32]byte) (*SinkWithdrawAppliedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}
	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}

	logs, sub, err := _Sink.contract.FilterLogs(opts, "WithdrawApplied", hashRule, txHashRule)
	if err != nil {
		return nil, err
	}
	return &SinkWithdrawAppliedIterator{contract: _Sink.contract, event: "WithdrawApplied", logs: logs, sub: sub}, nil
}

// WatchWithdrawApplied is a free log subscription operation binding the contract event 0x7f508fd15756f38a4426383f9ef243dfccdfdff0a528e755b113ef8bef2c5c2e.
//
// Solidity: event WithdrawApplied(hash indexed bytes32, txHash indexed bytes32, amount uint256, fee uint256, recipient address, category uint256, lastConfirmed address)
func (_Sink *SinkFilterer) WatchWithdrawApplied(opts *bind.WatchOpts, sink chan<- *SinkWithdrawApplied, hash [][32]byte, txHash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}
	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}

	logs, sub, err := _Sink.contract.WatchLogs(opts, "WithdrawApplied", hashRule, txHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SinkWithdrawApplied)
				if err := _Sink.contract.UnpackLog(event, "WithdrawApplied", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
