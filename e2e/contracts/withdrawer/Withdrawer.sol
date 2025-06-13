// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

struct Context {
    bytes origin;
    address sender;
    uint256 chainID;
}

interface IMRC20 {
    function approve(address spender, uint256 amount) external returns (bool);
    function withdraw(bytes memory to, uint256 amount) external returns (bool);
}

// Withdrawer is a simple contract performing a withdraw of deposited MRC20
// The amount to withdraw can be set during the contract deployment, it also to tests some edge cases like withdrawing BTC dust amount
contract Withdrawer {
    uint256 immutable public withdrawAmount;
    
    constructor(uint256 _withdrawAmount) {
        withdrawAmount = _withdrawAmount;
    }

    // perform a withdraw on cross chain call
    function onCrossChainCall(Context calldata context, address mrc20, uint256, bytes calldata) external {
        // perform withdrawal with the target token
        IMRC20(mrc20).approve(address(mrc20), type(uint256).max);
        IMRC20(mrc20).withdraw(context.origin, withdrawAmount);
    }

    // perform a withdraw on cross chain call, v2
    function onCall(Context calldata context, address mrc20, uint256, bytes calldata) external {
        // perform withdrawal with the target token
        IMRC20(mrc20).approve(address(mrc20), type(uint256).max);
        IMRC20(mrc20).withdraw(context.origin, withdrawAmount);
    }
}