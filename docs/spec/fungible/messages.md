# Messages

## MsgDeploySystemContracts

DeploySystemContracts deploy new instances of the system contracts

Authorized: admin policy group 2.

```proto
message MsgDeploySystemContracts {
	string creator = 1;
}
```

## MsgDeployFungibleCoinMRC20

DeployFungibleCoinMRC20 deploys a fungible coin from a connected chains as a MRC20 on MuseChain.

If this is a gas coin, the following happens:

* MRC20 contract for the coin is deployed
* contract address of MRC20 is set as a token address in the system
contract
* MUSE tokens are minted and deposited into the module account
* setGasMusePool is called on the system contract to add the information
about the pool to the system contract
* addLiquidityETH is called to add liquidity to the pool

If this is a non-gas coin, the following happens:

* MRC20 contract for the coin is deployed
* The coin is added to the list of foreign coins in the module's state

Authorized: admin policy group 2.

```proto
message MsgDeployFungibleCoinMRC20 {
	string creator = 1;
	string ERC20 = 2;
	int64 foreign_chain_id = 3;
	uint32 decimals = 4;
	string name = 5;
	string symbol = 6;
	pkg.coin.CoinType coin_type = 7;
	int64 gas_limit = 8;
	string liquidity_cap = 9;
}
```

## MsgRemoveForeignCoin

RemoveForeignCoin removes a coin from the list of foreign coins in the
module's state.

Authorized: admin policy group 2.

```proto
message MsgRemoveForeignCoin {
	string creator = 1;
	string mrc20_address = 2;
}
```

## MsgUpdateSystemContract

UpdateSystemContract updates the system contract

```proto
message MsgUpdateSystemContract {
	string creator = 1;
	string new_system_contract_address = 2;
}
```

## MsgUpdateContractBytecode

UpdateContractBytecode updates the bytecode of a contract from the bytecode
of an existing contract Only a MRC20 contract or the WMuse connector contract
can be updated IMPORTANT: the new contract bytecode must have the same
storage layout as the old contract bytecode the new contract can add new
variable but cannot remove any existing variable

Authozied: admin policy group 2

```proto
message MsgUpdateContractBytecode {
	string creator = 1;
	string contract_address = 2;
	string new_code_hash = 3;
}
```

## MsgUpdateMRC20WithdrawFee

UpdateMRC20WithdrawFee updates the withdraw fee and gas limit of a mrc20 token

```proto
message MsgUpdateMRC20WithdrawFee {
	string creator = 1;
	string mrc20_address = 2;
	string new_withdraw_fee = 6;
	string new_gas_limit = 7;
}
```

## MsgUpdateMRC20LiquidityCap

UpdateMRC20LiquidityCap updates the liquidity cap for a MRC20 token.

Authorized: admin policy group 2.

```proto
message MsgUpdateMRC20LiquidityCap {
	string creator = 1;
	string mrc20_address = 2;
	string liquidity_cap = 3;
}
```

## MsgPauseMRC20

PauseMRC20 pauses a list of MRC20 tokens
Authorized: admin policy group groupEmergency.

```proto
message MsgPauseMRC20 {
	string creator = 1;
	string mrc20_addresses = 2;
}
```

## MsgUnpauseMRC20

UnpauseMRC20 unpauses the MRC20 token
Authorized: admin policy group groupOperational.

```proto
message MsgUnpauseMRC20 {
	string creator = 1;
	string mrc20_addresses = 2;
}
```

## MsgUpdateGatewayContract

UpdateGatewayContract updates the mevm gateway contract used by the MuseChain protocol to read inbounds and process outbounds

```proto
message MsgUpdateGatewayContract {
	string creator = 1;
	string new_gateway_contract_address = 2;
}
```

## MsgUpdateMRC20Name

UpdateMRC20Name updates the name and/or the symbol of a mrc20 token

```proto
message MsgUpdateMRC20Name {
	string creator = 1;
	string mrc20_address = 2;
	string name = 3;
	string symbol = 4;
}
```

