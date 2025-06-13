//go:generate sh -c "solc GatewayMEVMCaller.sol --combined-json abi,bin | jq '.contracts.\"GatewayMEVMCaller.sol:GatewayMEVMCaller\"'  > GatewayMEVMCaller.json"
//go:generate sh -c "cat GatewayMEVMCaller.json | jq .abi > GatewayMEVMCaller.abi"
//go:generate sh -c "cat GatewayMEVMCaller.json | jq .bin  | tr -d '\"'  > GatewayMEVMCaller.bin"
//go:generate sh -c "abigen --abi GatewayMEVMCaller.abi --bin GatewayMEVMCaller.bin  --pkg gatewaymevmcaller --type GatewayMEVMCaller --out GatewayMEVMCaller.go"

package gatewaymevmcaller

var _ GatewayMEVMCaller
