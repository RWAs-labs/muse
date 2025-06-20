// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

struct Context {
    bytes origin;
    address sender;
    uint256 chainID;
}

interface mContract {
    function onCrossChainCall(
        Context calldata context,
        address mrc20,
        uint256 amount,
        bytes calldata message
    ) external;
}

contract ContextApp is mContract {

    constructor() {}
    event ContextData(bytes origin, address sender, uint256 chainID, address msgSender);

    // Call this function to perform a cross-chain swap
    function onCrossChainCall(Context calldata context, address mrc20, uint256 amount, bytes calldata message) external override {

        emit ContextData(context.origin, context.sender, context.chainID, msg.sender);
    }
}