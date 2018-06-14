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

pragma solidity ^0.4.10;

import "./oracle.sol";

// 本合约将记录审批流程签名哈希值以及确认提现申请
// 正常流程如下：
//  1. 调用 addHash 将审批流程签名哈希上链；
//  2. 由授权APP触发调用 enable 方法来确认本次哈希有效；
//  3. 提交提现申请 调用 approve 方法。
// 每次调用都需要经过2/3个系统签发者确认，才认可本次提交的数据正确性。
// 提现申请只能在审批流哈希上链之后，并且经过了授权APP的确认，才可以进行。
// 异常情况下，假如原审批流程其中一个环节发生了丢失密钥导致无法签名，可以启用
// 授权APP进行禁用该流程，再由原过程产生一个新流程。
// 禁用流程如下：
//  1. 调用 disable 方法。
// 提现申请需要传入本次提现的交易哈希和交易数额。由于需要经过2/3个系统签发者的确认，故需要每次都确认数字要对得上。
// 如果遇到其中一个有对不上的情况，则本次交易失败。
// This contract will save signature hash of approval flow and confirm withdrawal request
// A normal flow looks like:
//  1. Invoke addHash method to upload hash value of approval flow to private chain;
//  2. Manager App confirm the hash by invoking 'enable' method
//  3. Submit withdraw request by invoking 'approve' method
// Every invoktion has to be confirmed by at least 2/3 signers, and then we can confirm correctness of submission
// Withdrawl cannot be performed until approval flow has been uploaded to private chain and confirmed by manager APP
// If there is anything happen on any of node which result in signature failure, we can disable current flow using manager APP, and start a new flow.
// Disable flow works like as follows:
//  1. Invoke disable method
//  We have to send hash and amount as input parameters. Request has to be confirmed by at least 2/3 of signers, so the parameters have to be matched every time.
//  request will fail if any of these above doesn't match
contract Sink {

    Oracle oracle;

    // 审批流的四个状态
    // ADD     审批流哈希上链
    // ENABLE  确认审批流
    // DISABLE 禁用审批流
    // APPLY   申请提现
    // {-- 审批流哈希上链 --- 授权者确认上链的审批流哈希 --- 申请提现 --}
    // 申请提现必须是审批流哈希确认之后，未确认的审批流是不能通过申请提现的。
    // Status of approval flow
    // ADD      upload request to private chain
    // ENABLE   confirm request
    // DISABLE  disable request
    // APPLY    withdraw request
    // Withdraw request will not be accepted until approval flow on private chain has been confirmed
    enum Stage { ADD, ENABLE, DISABLE, APPLY }

    struct Counter {
        // msg.sender => bool
        mapping (address => bool) checked;
        address[] signers;
        // stage counter;
        uint count;

        // for apply withdraw
        uint amount;
        uint category;
        bytes32 txHash;
    }

    struct Journal {
        bytes32 hash;
        bool enabled;

        bool setuped;
        bool executed;

        // Stage => Counter
        mapping (uint8 => Counter) stages;
        // 申请过的相同交易不能再次重复申请,必须确保每次的体现申请交易号必须不一样
        // cannot have duplicated transactions, make sure transaction id is unique
        mapping (bytes32 => bool) trans;
    }

    // sign flow hash => journal index
    mapping (bytes32 => uint) ids;
    Journal[] journals;
    bool init;

    /* 系统签发事件 */
    event SignflowAdded(bytes32 hash, address lastConfirmed);
    event SignflowEnabled(bytes32 hash, address lastConfirmed);
    event SignflowDisabled(bytes32 hash, address lastConfirmed);
    // lastConfirmed 最后一个确认上链信息的signer地址
    event WithdrawApplied(bytes32 indexed hash,bytes32 indexed txHash, uint amount, uint fee, address recipient, uint category, address lastConfirmed);


    modifier onlySigner {
        require(oracle.isSigner(msg.sender));
        _;
    }

    // 本合约只能是n个节点中任意一个系统签名帐号来创建
    // Contract has to be created by any signature account in node
    function Sink(Oracle ref) public {
        oracle = ref;
        stageVerify(Stage.ADD, 0, 0, 0, 0, 0, 0);
        init = true;
    }

    function addHash(bytes32 hash) public onlySigner returns(bool) {
        return stageVerify(Stage.ADD, hash, 0, 0, 0, 0, 0);
    }

    function enable(bytes32 hash) public onlySigner returns(bool) {
        return stageVerify(Stage.ENABLE, hash, 0, 0, 0, 0, 0);
    }

    function disable(bytes32 hash) public onlySigner returns(bool) {
        return stageVerify(Stage.DISABLE, hash, 0, 0, 0, 0, 0);
    }

    // txHash 用于确定某一次的交易
    // category 转账的类别
    // txHash, used for confirm one transaction
    // category, transaction category
    function approve(bytes32 txHash, uint amount, uint fee, address recipient, bytes32 hash, uint category) public onlySigner returns(bool) {
        return stageVerify(Stage.APPLY, hash, txHash, amount, fee, recipient, category);
    }

    // 校验审批流程哈希是否有效,只记录已经上链的审批流程哈希，不确定未经过2/3系统签发者
    // Verify hash of approval flow is valid or not, only record those are on the private chain
    function available(bytes32 hash) constant public returns(bytes32, bool) {
        uint id = ids[hash];
        if (id == 0) {
            return (0x0, false);
        }

        return (journals[id].hash, journals[id].enabled);
    }

    // 只记录交易成功的交易哈希
    // Only record successful transaction hash
    function txExists(bytes32 hash, bytes32 txHash) constant public returns(bool) {
        uint id = ids[hash];
        if (id == 0) {
            return false;
        }

        Journal storage journal = journals[id];
        return journal.trans[txHash];
    }

    // 只有原始oracle中的boss才能更改oracle自身
    // 同时，新的oracle必须也要是boss才能生成
    // Only boss in original oracle can update itself
    // Only boss can create new oracle
    function changeOracle(address newOracle) public {
        require(oracle.boss() == msg.sender);
        oracle = Oracle(newOracle);
    }

    function stageVerify(Stage stage, bytes32 hash, bytes32 txHash, uint amount, uint fee, address recipient, uint category) internal returns(bool) {
        uint id = ids[hash];
        if (id == 0) {
            if (stage != Stage.ADD) {
                return false;
            }

            if (hash == 0 && init) {
                return false;
            }

            ids[hash] = journals.length;
            id = journals.length++;
        }

        Journal storage journal = journals[id];

        if (stage == Stage.ADD) {
            if (!journal.setuped) {
                journal.hash = hash;
                journal.setuped = true;
                if (!mark(journal, stage, txHash, amount, category)) {
                    return false;
                }

                // only one signer
                if (totalChecked(journal, stage) >= marginOfVotes()) {
                    SignflowAdded(hash, msg.sender);
                    journal.executed = true;
                }

                return true;
            }

            if (journal.executed || !mark(journal, stage, txHash, amount, category)) {
                return false;
            }
        }

        var (mustReturn, flag) = checkStageConditions(journal, stage, txHash, amount, category);
        if (mustReturn) {
            return flag;
        }

        if (stage == Stage.ADD) {
            SignflowAdded(hash, msg.sender);
            journal.executed = true;
        } else if (stage == Stage.ENABLE) {
            journal.enabled = true;
            reset(journal, stage);
            SignflowEnabled(hash, msg.sender);
        } else if (stage == Stage.DISABLE) {
            journal.enabled = false;
            reset(journal, stage);
            SignflowDisabled(hash, msg.sender);
        } else if (stage == Stage.APPLY) {
            journal.trans[txHash] = true;
            reset(journal, stage);
            WithdrawApplied(hash, txHash, amount, fee, recipient, category, msg.sender);
        } else {
            return false;
        }

        return true;
    }

    function checkStageConditions(Journal storage journal, Stage stage, bytes32 txHash, uint amount, uint category) internal returns(bool, bool) {
        if (stage == Stage.APPLY && journal.trans[txHash]) {
            return (true, false);
        }

        if (stage != Stage.ADD && (!journal.setuped || !journal.executed) ) {
            return (true, false);
        }

        if (stage != Stage.ADD && !mark(journal, stage, txHash, amount, category)) {
            return (true, false);
        }

        if (totalChecked(journal, stage) < marginOfVotes()) {
            return (true, true);
        }

        return (false, true);
    }

    function mark(Journal storage journal, Stage stage, bytes32 txHash, uint amount, uint category) internal returns(bool) {
        Counter storage counter = journal.stages[uint8(stage)];
        if ((stage == Stage.ENABLE && journal.enabled) || (stage == Stage.DISABLE && !journal.enabled) || (stage == Stage.APPLY && !journal.enabled)) {
            return false;
        }

        if (isChecked(journal, stage)) {
            return false;
        }

        if (stage == Stage.APPLY && counter.count == 0) {
            counter.amount = amount;
            counter.txHash = txHash;
            counter.category = category;
        } else if (stage == Stage.APPLY && counter.count > 0) {
            require(counter.amount == amount && counter.txHash == txHash && counter.category == category);
        }

        counter.checked[msg.sender] = true;
        counter.count++;

        if (stage != Stage.ADD) {
            uint len = counter.signers.length++;
            counter.signers[len] = msg.sender;
        }
        return true;
    }

    function reset(Journal storage journal, Stage stage) internal {
        Counter storage counter = journal.stages[uint8(stage)];

        uint i;
        uint len = counter.signers.length;
        for (i = 0; i < len; i++) {
            delete counter.checked[counter.signers[i]];
        }

        counter.count = 0;
        counter.signers.length = 0;
        counter.amount = 0;
        counter.txHash = 0x0;
    }

    function isChecked(Journal storage journal, Stage stage) constant internal returns(bool) {
        Counter storage counter = journal.stages[uint8(stage)];
        return counter.checked[msg.sender];
    }

    function totalChecked(Journal storage journal, Stage stage) constant internal returns(uint) {
        Counter storage counter = journal.stages[uint8(stage)];
        return counter.count;
    }

    // 服务的确认边界
    // Confirm border of services
    function marginOfVotes() constant internal returns(uint data) {
        uint totalNodes = oracle.totalEnabledNodes();
        // 只有一个节点或两个节点时
        // Only one or two nodes
        if (totalNodes == 1 || totalNodes == 2) {
            return totalNodes;
        }

        uint mod = totalNodes % 2;
        if (mod != 0) {
            // 奇数时，满足 2n + 1
            // Odd, satisfy 2n+1
            data = totalNodes/2 + mod;
        } else {
            data = totalNodes/2 + 1;
        }

        return data;
    }

}