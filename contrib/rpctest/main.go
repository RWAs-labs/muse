package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	museeth "github.com/RWAs-labs/protocol-contracts/pkg/muse.eth.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// #nosec G101 - used for testing only
	MuseEthPriv           = "9D00E4D7A8A14384E01CD90B83745BCA847A66AD8797A9904A200C28C2648E64"
	SystemContractAddress = "0x91d18e54DAf4F677cB28167158d6dd21F6aB3921"
)

type Request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *Error          `json:"error"`
	ID     int             `json:"id"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <blocknum>\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("Start testing the mEVM ETH JSON-RPC for all txs...\n")
	fmt.Printf("Test1: simple gas voter tx\n")

	bn, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	if bn < 0 {
		panic("Block number must be non-negative")
	}
	// #nosec G115 check as positive
	bnUint64 := uint64(bn)

	if false {
		// THIS WOULD NOT WORK: see https://github.com/RWAs-labs/muse-node/pull/445
		// USE RAW JSON-RPC INSTEAD
		mevmClient, err := ethclient.Dial("http://localhost:8545")
		if err != nil {
			panic(err)
		}

		block, err := mevmClient.BlockByNumber(context.Background(), big.NewInt(bn))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Block number: %d, num of txs %d (should be 1)\n", block.Number(), len(block.Transactions()))
	}

	client := &EthClient{
		Endpoint:   "http://localhost:8545",
		HTTPClient: &http.Client{},
	}
	resp := client.EthGetBlockByNumber(bnUint64, false)
	var jsonObject map[string]interface{}
	if resp.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", resp.Error.Message, resp.Error.Code)
		panic(resp.Error.Message)
	}
	err = json.Unmarshal(resp.Result, &jsonObject)
	if err != nil {
		panic(err)
	}

	txs, ok := jsonObject["transactions"].([]interface{})
	if !ok || len(txs) != 1 {
		panic("Wrong number of txs")
	}
	txhash, ok := txs[0].(string)
	if !ok {
		panic("Wrong tx type")
	}
	fmt.Printf("Tx hash: %s\n", txhash)
	tx := client.EthGetTransactionReceipt(txhash)
	if tx.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", tx.Error.Message, tx.Error.Code)
		panic(tx.Error.Message)
	}

	// tx receipt can be queried by ethclient queries.
	mevmClient, err := ethclient.Dial(client.Endpoint)
	if err != nil {
		panic(err)
	}
	receipt, err := mevmClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(txhash))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Receipt status: %+v\n", receipt.Status)

	// HeaderByHash works; BlockByHash does not work;
	// main offending RPC is the transaction type; we have custom type id 56
	// which is not recognized by the go-ethereum client.
	blockHeader, err := mevmClient.HeaderByNumber(context.Background(), big.NewInt(bn))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Block header TxHash: %+v\n", blockHeader.TxHash)

	chainid, err := mevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	museEthPrivKey, err := crypto.HexToECDSA(MuseEthPriv)
	if err != nil {
		panic(err)
	}
	mevmAuth, err := bind.NewKeyedTransactorWithChainID(museEthPrivKey, chainid)
	if err != nil {
		panic(err)
	}
	museContractAddress, tx2, museContract, err := museeth.DeployMuseEth(
		mevmAuth,
		mevmClient,
		mevmAuth.From,
		big.NewInt(2_100_000_000),
	)
	_, _ = museContractAddress, museContract
	if err != nil {
		panic(err)
	}
	//time.Sleep(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	receipt = MustWaitForReceipt(ctx, mevmClient, tx2.Hash())
	receipt, err = mevmClient.TransactionReceipt(context.Background(), tx2.Hash())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deploy EthMuse Contract Receipt: %+v\n", receipt)
	receipt2 := client.EthGetTransactionReceipt(tx2.Hash().Hex())
	if receipt2.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", receipt2.Error.Message, receipt2.Error.Code)
		panic(tx.Error.Message)
	}
	jsonObject = make(map[string]interface{})
	err = json.Unmarshal(receipt2.Result, &jsonObject)
	if err != nil {
		panic(err)
	}
	prettyJSON, err := json.MarshalIndent(jsonObject, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %s\n", string(prettyJSON))

	fmt.Printf("MuseEth Contract Address: %s\n", museContractAddress.Hex())
	if museContractAddress != receipt.ContractAddress {
		panic(fmt.Sprintf("Contract address mismatch: wanted %s, got %s", museContractAddress, receipt.ContractAddress))
	}

	{
		//  test of filter event of MuseEth, a typical ERC20 token contract
		tx, err := museContract.Approve(mevmAuth, ethcommon.HexToAddress(SystemContractAddress), big.NewInt(100))
		if err != nil {
			panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		receipt = MustWaitForReceipt(ctx, mevmClient, tx.Hash())
		fmt.Printf("Approve Receipt status: %+v\n", receipt.Status)
		approveIter, err := museContract.FilterApproval(&bind.FilterOpts{
			Start: 0,
			End:   nil,
		}, nil, nil)
		if err != nil {
			panic(err)
		}
		for approveIter.Next() {
			event := approveIter.Event
			fmt.Printf(
				"Approval event: owner %s, spender %s, amount %s\n",
				event.Owner.Hex(),
				event.Spender.Hex(),
				event.Value.String(),
			)
			fmt.Printf("  raw log txhash: %s\n", event.Raw.TxHash.Hex())
			fmt.Printf("  tx hash: %s\n", tx.Hash().Hex())
			if event.Raw.TxHash != tx.Hash() {
				panic("Tx hash mismatch")
			}
		}
	}

	// test getLogs
	sys, err := systemcontract.NewSystemContract(ethcommon.HexToAddress(SystemContractAddress), mevmClient)
	if err != nil {
		panic(err)
	}
	gas, err := sys.GasPriceByChainId(&bind.CallOpts{}, big.NewInt(1337))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Gas price: %s\n", gas.String())
	toBlock := bnUint64
	gasPriceIter, err := sys.FilterSetGasPrice(&bind.FilterOpts{
		Start: bnUint64,
		End:   &toBlock,
	})
	if err != nil {
		panic(err)
	}

	for gasPriceIter.Next() {
		event := gasPriceIter.Event
		fmt.Printf("Set gas price event: arg0 %d, arg1 %d\n", event.Arg0, event.Arg1)
		fmt.Printf("  raw event: %+v\n", event.Raw)

		if event.Raw.TxHash.Big().Cmp(big.NewInt(0)) == 0 {
			panic("TxHash is zero")
		}
		fmt.Printf("  raw log txhash: %s\n", event.Raw.TxHash.Hex())
		fmt.Printf("  querying the receipt...\n")
		receipt = MustWaitForReceipt(context.Background(), mevmClient, event.Raw.TxHash)
		fmt.Printf("  receipt: hash: %s\n", receipt.TxHash)
		fmt.Printf("    receipt status: %d\n", receipt.Status)
		fmt.Printf("    receipt gas used: %d\n", receipt.GasUsed)
		fmt.Printf("    receipt cumulative gas used: %d\n", receipt.CumulativeGasUsed)
		fmt.Printf("    receipt logs: %+v\n", receipt.Logs)
		fmt.Printf("    receipt block hash: %s\n", receipt.BlockHash.Hex())
		fmt.Printf("    receipt block number: %d\n", receipt.BlockNumber)
		fmt.Printf("    receipt transaction index: %d\n", receipt.TransactionIndex)
		fmt.Printf("    receipt contract address: %s\n", receipt.ContractAddress.Hex())
	}
}

type EthClient struct {
	Endpoint   string
	HTTPClient *http.Client
}

func (c *EthClient) EthGetBlockByNumber(blockNum uint64, verbose bool) *Response {
	client := c.HTTPClient
	hexBlockNum := fmt.Sprintf("0x%x", blockNum)
	req := &Request{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			hexBlockNum,
			verbose,
		},
		ID: 1,
	}

	// Encode the request to JSON
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		panic(err)
	}
	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", c.Endpoint, buf)
	if err != nil {
		panic(err)
	}
	// Set the content type header
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	// #nosec G107 - defer close
	defer resp.Body.Close()
	// Decode the response from JSON
	var rpcResp Response
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		panic(err)
	}

	return &rpcResp
}

func (c *EthClient) EthGetTransactionReceipt(txhash string) *Response {
	client := c.HTTPClient
	req := &Request{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionReceipt",
		Params: []interface{}{
			txhash,
		},
		ID: 1,
	}

	// Encode the request to JSON
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		panic(err)
	}
	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", c.Endpoint, buf)
	if err != nil {
		panic(err)
	}
	// Set the content type header
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	// #nosec G107 - defer close
	defer resp.Body.Close()
	// Decode the response from JSON
	var rpcResp Response
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		panic(err)
	}

	return &rpcResp
}

func (c *EthClient) EthGetLogs() {
	//client := c.HTTPClient
	//req := &Request{
	//	Jsonrpc: "2.0",
	//	Method:  "eth_getTransactionReceipt",
	//	Params: []interface{}{
	//		txhash,
	//	},
	//	ID: 1,
	//}
	//
	//// Encode the request to JSON
	//buf := &bytes.Buffer{}
	//err := json.NewEncoder(buf).Encode(req)
	//if err != nil {
	//	panic(err)
	//}
	//// Create a new HTTP request
	//httpReq, err := http.NewRequest("POST", c.Endpoint, buf)
	//if err != nil {
	//	panic(err)
	//}
	//// Set the content type header
	//httpReq.Header.Set("Content-Type", "application/json")
	//
	//// Send the HTTP request
	//resp, err := client.Do(httpReq)
	//if err != nil {
	//	panic(err)
	//}
	//// #nosec G107 - defer close
	//defer resp.Body.Close()
	//// Decode the response from JSON
	//var rpcResp Response
	//err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	//if err != nil {
	//	panic(err)
	//}
	//
	//return &rpcResp
}

func MustWaitForReceipt(ctx context.Context, client *ethclient.Client, txhash ethcommon.Hash) *types.Receipt {
	for {
		select {
		case <-ctx.Done():
			panic("timeout waiting for transaction receipt")
		default:
			receipt, err := client.TransactionReceipt(context.Background(), txhash)
			if err == nil {
				return receipt
			}
			time.Sleep(1 * time.Second)
		}
	}
}
