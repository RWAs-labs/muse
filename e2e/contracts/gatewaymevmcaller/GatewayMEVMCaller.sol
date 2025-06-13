// SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

struct CallOptions {
    uint256 gasLimit;
    bool isArbitraryCall;
}

struct RevertOptions {
    address revertAddress;
    bool callOnRevert;
    address abortAddress;
    bytes revertMessage;
    uint256 onRevertGasLimit;
}

interface IGatewayMEVM {
    function call(
        bytes memory receiver,
        address mrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    )
        external;

    function withdrawAndCall(
        bytes memory receiver,
        uint256 amount,
        uint256 chainId,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    )
        external;

    function withdrawAndCall(
        bytes memory receiver,
        uint256 amount,
        address mrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    )
        external;
}

interface IMRC20 {
    function approve(address spender, uint256 amount) external returns (bool);
}

interface WMUSE {
    function deposit() external payable;
    function approve(address guy, uint256 wad) external returns (bool);
}

contract GatewayMEVMCaller {
    IGatewayMEVM private gatewayMEVM;
    WMUSE wmuse;
    constructor(address gatewayMEVMAddress, address wmuseAddress) {
        gatewayMEVM = IGatewayMEVM(gatewayMEVMAddress);
        wmuse = WMUSE(wmuseAddress);
    }

    function callGatewayMEVM(
        bytes memory receiver,
        address mrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    ) external {
        IMRC20(mrc20).approve(address(gatewayMEVM), 100000000000000000);
        gatewayMEVM.call(receiver, mrc20, message, callOptions, revertOptions);
    }

    function withdrawAndCallGatewayMEVM(
        bytes memory receiver,
        uint256 amount,
        uint256 chainId,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    ) external {
        wmuse.approve(address(gatewayMEVM), amount);
        gatewayMEVM.withdrawAndCall(receiver, amount, chainId, message, callOptions, revertOptions);
    }

    function withdrawAndCallGatewayMEVM(
        bytes memory receiver,
        uint256 amount,
        address mrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    ) external {
        IMRC20(mrc20).approve(address(gatewayMEVM), 100000000000000000);
        gatewayMEVM.withdrawAndCall(receiver, amount, mrc20, message, callOptions, revertOptions);
    }

    function depositWMUSE() external payable {
        wmuse.deposit{value: msg.value}();
    }
}