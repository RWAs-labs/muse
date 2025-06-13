//go:generate sh -c "solc --evm-version paris MEVMSwapApp.sol --combined-json abi,bin --allow-paths .. | jq '.contracts.\"MEVMSwapApp.sol:MEVMSwapApp\"'  > MEVMSwapApp.json"
//go:generate sh -c "cat MEVMSwapApp.json | jq .abi > MEVMSwapApp.abi"
//go:generate sh -c "cat MEVMSwapApp.json | jq .bin  | tr -d '\"'  > MEVMSwapApp.bin"
//go:generate sh -c "abigen --abi MEVMSwapApp.abi --bin MEVMSwapApp.bin --pkg mevmswap --type MEVMSwapApp --out MEVMSwapApp.go"

package mevmswap
