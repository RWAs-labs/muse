# `musee2e`

### Basics

`musee2e` is a CLI tool allowing to quickly test MuseChain functionality.

`musee2e` can communicate with MuseChain, EVM and Bitcoin network to test and track the full workflow of a cross-chain transaction.

### Getting Started

`musee2e` can be installed with the command:

```go
make install-musee2e

musee2e -h
```

### Config

The command takes a config file as input containing RPC nodes to connect to, hotkey wallet information for interaction with networks, and addresses of the deployed contracts to be used.

This is an example of config for interaction with Athens3:

```go
muse_chain_id: "athens_7001-1"
default_account:
  evm_address: "<your evm address>"
  private_key: "<your evm private key>"
rpcs:
  mevm: "<mevm (MuseChain) url>, generally using port 8545"
  evm: "<evm url>, generally using port 8545"
  bitcoin:
    host: "<bitcoin rpc url>"
    user: "<bitcoin user>"
    pass: "<bitcoin pass>"
    http_post_mode: true
    disable_tls: true
    params: "<mainnet|testnet3|regnet>"
  musecore_grpc: "<MuseChain grpc url>, generally using port 9090"
  musecore_rpc: "<MuseChain grpc url>, generally using port 26657"
contracts:
  mevm:
    system_contract: "0xEdf1c3275d13489aCdC6cD6eD246E72458B8795B"
    eth_mrc20: "0x13A0c5930C028511Dc02665E7285134B6d11A5f4"
    erc20_mrc20: "0x0cbe0dF132a6c6B4a2974Fa1b7Fb953CF0Cc798a"
    btc_mrc20: "0x65a45c57636f9BcCeD4fe193A602008578BcA90b"
    uniswap_factory: "0x9fd96203f7b22bCF72d9DCb40ff98302376cE09c"
    uniswap_router: "0x2ca7d64A7EFE2D62A725E2B35Cf7230D6677FfEe"
  evm:
    muse_eth: "0x0000c304d2934c00db1d51995b9f6996affd17c0"
    connector_eth: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea"
    custody: "0x000047f11c6e42293f433c82473532e869ce4ec5"
    erc20: "0x07865c6e87b9f70255377e024ace6630c1eaa37f"
```

### Bitcoin setup
Interaction with the Bitcoin node will require setting up a specific node tracking the address.

If an error occurs mention that wallets are not loaded, this command can be run to set it up:

```go
musee2e setup-bitcoin [config]
```

Note: This command is intended for local testing only and won't work on live networks, where a node with an enabled wallet should already be set up.

### Commands

Show the balances of the accounts used on the different networks:

```go
musee2e balances [config]
```

Show the Bitcoin address (the address is derived from the Ethereum private key, this address must therefore be found to perform the Bitcoin test):

```go
musee2e bitcoin-address [config]
```

The list of tests to be run can be found by running following command:

```go
musee2e list-tests
```

Run tests:

Tests can be run using `run` command. Arguments of this command should be in `testName:testArgs` format, where `testArgs` is optional comma-separated list of test arguments. If test arguments are not provided, default arguments will be used. Example of `run` command:

```go
musee2e run muse_deposit:2000000000000000000 eth_deposit:2000000000000000000 erc20_deposit:200000 --config cmd/musee2e/config/local.yml
```

Since cctxs might take a longer time to be processed on live networks, it is highly recommended to use `--verbose` flag to see the current status of the cctx workflow.

### Testing a gas MRC20 from an EVM chain

Testing a gas token requires the following values to be defined in the config:

```go
muse_chain_id
default_account:
  evm_address
  private_key
rpcs:
  mevm
  evm
  musecore_grpc
  musecore_rpc
contracts:
  mevm:
    system_contract
    eth_mrc20
    uniswap_factory
    uniswap_router
  evm:
    muse_eth
    connector_eth
    custody: "0x000047f11c6e42293f433c82473532e869ce4ec5"
```

One of the tests can be commented out in case only a deposit or a withdrawal is to be tested.
Testing an ERC20 MRC20 from an EVM chain

Testing MRC20 requires the same config as for the gas tokens, but must include the `erc20` field that contains the address of the ERC20 on the evm chain and `erc20_mrc20` on MuseChain.

```go
muse_chain_id
accounts:
  evm_address
  evm_priv_key
rpcs:
  mevm
  evm
  musecore_grpc
  musecore_rpc
contracts:
  mevm:
    erc20_mrc20
  evm:
		erc20
```

### Testing a MRC20 from a Bitcoin chain

The following values must be set in the config in order to test Bitcoin functionality

```go
muse_chain_id
accounts:
  evm_address
  evm_priv_key
rpcs:
  mevm
  bitcoin:
    host
    user
    pass
    http_post_mode
    disable_tls
    params
  musecore_grpc
  musecore_rpc
contracts:
  mevm:
    system_contract
    btc_mrc20
    uniswap_factory
    uniswap_router
```

### TODO: message passing
