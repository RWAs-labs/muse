// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @title IBank Interface for Cross-chain Token Deposits and Withdrawals
/// @notice This interface defines the functions for depositing MRC20 tokens and withdrawing Cosmos tokens,
///         as well as querying the balance of Cosmos tokens corresponding to a given MRC20 token.
/// @dev This contract interacts with a precompiled contract at a fixed address.

/// @dev The IBank contract's precompiled address.
address constant IBANK_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000067; // Address 103

/// @dev The IBank contract instance using the precompiled address.
IBank constant IBANK_CONTRACT = IBank(IBANK_PRECOMPILE_ADDRESS);

/// @dev Interface for the IBank contract.
interface IBank {
    /// @notice Deposit event is emitted when deposit function is called.
    /// @param mrc20_depositor Depositor EVM address.
    /// @param mrc20_token MRC20 address deposited.
    /// @param cosmos_token Cosmos token denomination the tokens were converted into.
    /// @param cosmos_address Cosmos address the tokens were deposited to.
    /// @param amount Amount deposited.
    event Deposit(
        address indexed mrc20_depositor,
        address indexed mrc20_token,
        string indexed cosmos_token,
        string cosmos_address,
        uint256 amount
    );

    /// @notice Withdraw event is emitted when withdraw function is called.
    /// @param mrc20_withdrawer Withdrawer EVM address.
    /// @param mrc20_token MRC20 address withdrawn.
    /// @param cosmos_token Cosmos token denomination the tokens were converted from.
    /// @param cosmos_address Cosmos address the tokens were withdrawn from.
    /// @param amount Amount withdrawn.
    event Withdraw(
        address indexed mrc20_withdrawer,
        address indexed mrc20_token,
        string indexed cosmos_token,
        string cosmos_address,
        uint256 amount
    );

    /// @notice Deposit a MRC20 token and mint the corresponding Cosmos token to the user's account.
    /// @param mrc20 The MRC20 token address to be deposited.
    /// @param amount The amount of MRC20 tokens to deposit.
    /// @return success Boolean indicating whether the deposit was successful.
    function deposit(
        address mrc20,
        uint256 amount
    ) external returns (bool success);

    /// @notice Withdraw Cosmos tokens and convert them back to the corresponding MRC20 token for the user.
    /// @param mrc20 The MRC20 token address for the corresponding Cosmos token.
    /// @param amount The amount of Cosmos tokens to withdraw.
    /// @return success Boolean indicating whether the withdrawal was successful.
    function withdraw(
        address mrc20,
        uint256 amount
    ) external returns (bool success);

    /// @notice Retrieve the Cosmos token balance corresponding to a specific MRC20 token for a given user.
    /// @param mrc20 The MRC20 cosmos token denomination to check the balance for.
    /// @param user The address of the user to retrieve the balance for.
    /// @return balance The balance of the Cosmos token for the specified MRC20 token and user.
    function balanceOf(
        address mrc20,
        address user
    ) external view returns (uint256 balance);
}
