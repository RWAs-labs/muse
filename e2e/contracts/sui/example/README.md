# SUI WithdrawAndCall with PTB Transactions

This document explains how the SUI `withdrawAndCall` functionality works using Programmable Transaction Blocks (PTB) in the MuseChain protocol.

## Overview

The `withdrawAndCall` operation in MuseChain allows users to withdraw tokens from MEVM to the Sui blockchain and simultaneously calls a `on_call` function in the `connected` module on the Sui side.

This is implemented as a single atomic transaction using Sui's Programmable Transaction Blocks (PTB).

## Transaction Flow

1. **User Initiates Withdrawal**: A user initiates a withdrawal from MEVM to Sui with a `on_call` payload.

2. **MEVM Processing**: The MEVM gateway processes the withdrawal request and prepares the transaction.

3. **PTB Construction**: A Programmable Transaction Block is constructed with the following steps:
   - **Withdraw**: The first command in the PTB is the `withdraw_impl` function call, which:
     - Verifies the withdrawal parameters
     - Withdraw and returns two coin objects: the main withdrawn coins and the gas budget coins
   - **Gas Budget Transfer**: The second command transfers the gas budget coins to the TSS address to cover transaction fees.
     - The gas budget is the SUI coin withdrawn from sui vault, together with withdrawn CCTX's coinType.
     - The gas budget needs to be forwarded to TSS address to cover the transaction fee.
   - **Connected Module Call**: The third command calls the `on_call` function in the connected module, passing:
     - The withdrawn coins
     - The call payload from the user
     - Any additional parameters required by the connected module

4. **Transaction Execution**: The entire PTB is executed atomically on the Sui blockchain.

## PTB Structure

The PTB for a `withdrawAndCall` transaction consists of three main commands:

```
PTB {
    // Command 0: Withdraw Implementation
    MoveCall {
        package: gateway_package_id,
        module: gateway_module,
        function: withdraw_impl,
        arguments: [
            gateway_object_ref,
            withdraw_cap_object_ref,
            coin_type,
            amount,
            nonce,
            gas_budget
        ]
    }
    
    // Command 1: Gas Budget Transfer
    TransferObjects {
        from: withdraw_impl_result[1], // Gas budget coins
        to: tss_address
    }
    
    // Command 2: Connected Module Call
    MoveCall {
        package: target_package_id,
        module: connected_module,
        function: on_call,
        arguments: [
            withdraw_impl_result[0], // Main withdrawn coins
            on_call_payload
        ]
    }
}
``` 