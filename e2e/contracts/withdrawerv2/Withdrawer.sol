// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface IMRC20 {
    function approve(address spender, uint256 amount) external;
    function transferFrom(address sender, address recipient, uint256 amount) external;
    function withdraw(bytes memory to, uint256 amount) external returns (bool);
    function withdrawGasFee() external returns (address, uint256);
}

// Sample contract for running withdraw on mEVM
contract Withdrawer {
    // Run n withdraws of amount on asset to custody
    function runWithdraws(
        bytes calldata recipient,
        IMRC20 asset,
        uint256 amount,
        uint256 count
    ) external {
        // transfer gas for the transactions and approve it in the mrc20
        (address gas, uint256 gasFee) = asset.withdrawGasFee();
        IMRC20(gas).transferFrom(msg.sender, address(this), gasFee * (count+1));
        IMRC20(gas).approve(address(asset), gasFee * (count+1));

        // perform the withdraws
        asset.transferFrom(msg.sender, address(this), amount * count);
        for (uint256 i = 0; i < count; i++) {
            asset.withdraw(recipient, amount);
        }
    }
}