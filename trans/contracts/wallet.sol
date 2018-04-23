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

import "./bank.sol";
import "./erc20.sol";

contract Wallet {

    event WalletS(address bank);

    Bank bank;
    address owner;

    function Wallet(Bank _bank) public {
        bank = _bank;
        owner = msg.sender;
        WalletS(bank);
    }

    function() payable {
        if (msg.value > 0) {
            bank.transfer(msg.value);
        }
    }

    function transferERC20(ERC20 token, address _to, uint _amount) {
        // 只有收款账号才能操作转账
        require(msg.sender == owner);
        if (!token.transfer(_to, _amount)) {
            revert();
        }
    }
}
