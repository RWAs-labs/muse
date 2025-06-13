// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "interfaces/IUniswapV2Router02.sol";

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

interface IMRC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool);
    function deposit(address to, uint256 amount) external returns (bool);
    function burn(address account, uint256 amount) external returns (bool);
    function withdraw(bytes memory to, uint256 amount) external returns (bool);
    function withdrawGasFee() external view returns (address, uint256);
    function PROTOCOL_FEE() external view returns (uint256);
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Deposit(bytes from, address indexed to, uint256 value);
    event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasFee, uint256 protocolFlatFee);
}

contract MEVMSwapApp is mContract {
    error InvalidSender();
    error LowAmount();

    uint256 constant private _DEADLINE = 1 << 64;
    address immutable public router02;
    address immutable public systemContract;
    
    constructor(address router02_, address systemContract_) {
        router02 = router02_;
        systemContract = systemContract_;
    }

    function encodeMemo(address targetMRC20,  bytes calldata recipient) pure external returns (bytes memory) {
        return abi.encodePacked(targetMRC20, recipient);
    }

    function decodeMemo(bytes calldata data) pure public returns (address, bytes memory) {
        bytes memory decodedBytes;
        uint256 size;
        address addr;

        size = data.length;
        addr = address(uint160(bytes20(data[0:20])));
        decodedBytes = data[20:];

        return (addr, decodedBytes);
    }

    // Perform a cross-chain swap
    function onCrossChainCall(Context calldata, address mrc20, uint256 amount, bytes calldata message) external override {
        if (msg.sender != systemContract) {
            revert InvalidSender();
        }

        address targetMRC20;
        bytes memory recipient;
        address[] memory path;

        (targetMRC20, recipient) = decodeMemo(message);
        path = new address[](2);
        path[0] = mrc20;
        path[1] = targetMRC20;

        // approve the usage of this token by router02
        IMRC20(mrc20).approve(address(router02), amount);

        // swap for target token
        uint256[] memory amounts = IUniswapV2Router02(router02).swapExactTokensForTokens(amount, 0, path, address(this), _DEADLINE);

        // perform withdrawal with the target token
        (address gasMRC20Addr,uint256 gasFee) = IMRC20(targetMRC20).withdrawGasFee();
        IMRC20(gasMRC20Addr).approve(address(targetMRC20), gasFee);
        IMRC20(targetMRC20).approve(address(targetMRC20), amounts[1]*10);
        IMRC20(targetMRC20).withdraw(recipient, amounts[1]);
    }

    // used with v2 contracts
    function onCall(Context calldata, address mrc20, uint256 amount, bytes calldata message) external {
        address targetMRC20;
        bytes memory recipient;
        address[] memory path;

        (targetMRC20, recipient) = decodeMemo(message);
        path = new address[](2);
        path[0] = mrc20;
        path[1] = targetMRC20;

        // approve the usage of this token by router02
        IMRC20(mrc20).approve(address(router02), amount);

        // swap for target token
        uint256[] memory amounts = IUniswapV2Router02(router02).swapExactTokensForTokens(amount, 0, path, address(this), _DEADLINE);

        // perform withdrawal with the target token
        (address gasMRC20Addr,uint256 gasFee) = IMRC20(targetMRC20).withdrawGasFee();
        IMRC20(gasMRC20Addr).approve(address(targetMRC20), gasFee);
        IMRC20(targetMRC20).approve(address(targetMRC20), amounts[1]*10);
        IMRC20(targetMRC20).withdraw(recipient, amounts[1]);
    }
}