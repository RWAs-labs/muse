//go:generate sh -c "solc --evm-version paris TestMuseConnectorMEVM.sol --combined-json abi,bin | jq '.contracts.\"TestMuseConnectorMEVM.sol:TestMuseConnectorMEVM\"'  > TestMuseConnectorMEVM.json"
//go:generate sh -c "cat TestMuseConnectorMEVM.json | jq .abi > TestMuseConnectorMEVM.abi"
//go:generate sh -c "cat TestMuseConnectorMEVM.json | jq .bin  | tr -d '\"'  > TestMuseConnectorMEVM.bin"
//go:generate sh -c "abigen --abi TestMuseConnectorMEVM.abi --bin TestMuseConnectorMEVM.bin --pkg testconnectormevm --type TestMuseConnectorMEVM --out TestMuseConnectorMEVM.go"

package testconnectormevm
