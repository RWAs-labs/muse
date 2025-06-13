//go:generate sh -c "solc --evm-version london TestMRC20.sol --combined-json abi,bin | jq '.contracts.\"TestMRC20.sol:TestMRC20\"'  > TestMRC20.json"
//go:generate sh -c "cat TestMRC20.json | jq .abi > TestMRC20.abi"
//go:generate sh -c "cat TestMRC20.json | jq .bin  | tr -d '\"'  > TestMRC20.bin"
//go:generate sh -c "abigen --abi TestMRC20.abi --bin TestMRC20.bin --pkg testmrc20 --type TestMRC20 --out TestMRC20.go"

package testmrc20
