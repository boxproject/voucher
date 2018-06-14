pragma solidity ^0.4.10;

import "./erc20.sol";

contract Bank {
    // Deposit eth充值事件; eth transfer event
    // account 充值的人; account transfer from
    // amount 充值数量; amount to transfer
    event Deposit(address account, uint256 amount);
    event Withdraw(address to, uint256 amount);
    event AllowFlow(bytes32 hash);
    event DisallowFlow(bytes32 hash);

    address operator;
    mapping (bytes32 => bool) signFlow;

    modifier onlyOperator {
        require(operator == msg.sender);
        _;
    }

    modifier onlySignflow(bytes32 sf) {
        require(signFlow[sf]);
        _;
    }

    function Bank() public {
        // 该帐号拥有超级超级权限; people who created the contract own it.
        operator = msg.sender;
    }

    // 任何人都可以向这个合约充值
    // 记录日志 谁充值的，充到哪个钱包，充了多少钱
    // Anyone can transfer fund to bank
    function() payable public {
        if (msg.value > 0) {
            Deposit(msg.sender, msg.value);
        }
    }

    // eth提现
    // eth withdrawl, only owner with valid signature can withdraw
    function withdraw(address to, uint amount, bytes32 sf) public onlyOperator onlySignflow(sf) {
        require(to != 0 && amount > 0 && this.balance >= amount);
        to.transfer(amount);
        Withdraw(to, amount);
    }

    function transferERC20(ERC20 token, address _to, uint _amount, bytes32 sf) public onlyOperator onlySignflow(sf) {
        require(_to != 0 && _amount > 0);
        // 正确性交由erc20合约来决定
        if (!token.transfer(_to, _amount)) {
            revert();
        }
    }

    function allow(bytes32 hash) public onlyOperator returns(bool) {
        if (signFlow[hash]) {
            return false;
        }

        signFlow[hash] = true;
        AllowFlow(hash);
        return true;
    }

    function disallow(bytes32 hash) public onlyOperator returns(bool) {
        if (signFlow[hash]) {
            delete signFlow[hash];
            return true;
        }

        DisallowFlow(hash);
        return false;
    }

    function available(bytes32 sf) constant public returns(bool) {
        return signFlow[sf];
    }
}