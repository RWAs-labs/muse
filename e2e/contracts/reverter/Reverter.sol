// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Sample contract for evm tests
contract Reverter {
    error Foo();

    struct mContext {
        bytes origin;
        address sender;
        uint256 chainID;
    }

    function onCrossChainCall(
        mContext calldata context,
        address mrc20,
        uint256 amount,
        bytes calldata message
    ) external {
        onCall(context, mrc20, amount, message);
    }

    function onCall(
        mContext calldata context,
        address mrc20,
        uint256 amount,
        bytes calldata message
    ) public {
        revert Foo();
    }
}