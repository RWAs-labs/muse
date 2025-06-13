---
sidebar_position: 1
---

# Overview

The `fungible` module facilitates the deployment of fungible tokens of connected
blockchains (called "foreign coins") on MuseChain.

Foreign coins are represented as MRC20 tokens on MuseChain.

When a foreign coin is deployed on MuseChain, a MRC20 contract is deployed, a
pool is created, liquidity is added to the pool, and the foreign coin is added
to the list of foreign coins in the module's state.

The module contains the logic for:

- Deploying a foreign coin on MuseChain
- Deploying a system contract, Uniswap and wrapped MUSE
- Depositing to and calling omnichain smart contracts on MuseChain from
  connected chains (`DepositMRC20AndCallContract` and `DepositMRC20`)

the module depends heavily on the
[protocol contracts](https://github.com/RWAs-labs/protocol-contracts).

## State

The `fungible` module keeps track of the following state:

- System contract address
- A list of foreign coins
