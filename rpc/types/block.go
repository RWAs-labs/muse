// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/RWAs-labs/ethermint/blob/main/LICENSE
package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	ethermint "github.com/RWAs-labs/ethermint/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

// BlockNumber represents decoding hex string to block values
type BlockNumber int64

const (
	EthPendingBlockNumber  = BlockNumber(-2)
	EthLatestBlockNumber   = BlockNumber(-1)
	EthEarliestBlockNumber = BlockNumber(0)
)

const (
	BlockParamEarliest  = "earliest"
	BlockParamLatest    = "latest"
	BlockParamFinalized = "finalized"
	BlockParamSafe      = "safe"
	BlockParamPending   = "pending"
)

// NewBlockNumber creates a new BlockNumber instance.
func NewBlockNumber(n *big.Int) BlockNumber {
	if !n.IsInt64() {
		// default to latest block if it overflows
		return EthLatestBlockNumber
	}

	return BlockNumber(n.Int64())
}

// ContextWithHeight wraps a context with the a gRPC block height header. If the provided height is
// 0, it will return an empty context and the gRPC query will use the latest block height for querying.
// Note that all metadata are processed and removed by tendermint layer, so it wont be accessible at gRPC server level.
func ContextWithHeight(height int64) context.Context {
	if height == 0 {
		return context.Background()
	}

	return metadata.AppendToOutgoingContext(
		context.Background(),
		grpctypes.GRPCBlockHeightHeader,
		fmt.Sprintf("%d", height),
	)
}

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "latest", "finalized", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case BlockParamEarliest:
		*bn = EthEarliestBlockNumber
		return nil
	case BlockParamLatest, BlockParamFinalized, BlockParamSafe:
		*bn = EthLatestBlockNumber
		return nil
	case BlockParamPending:
		*bn = EthPendingBlockNumber
		return nil
	}

	blckNum, err := hexutil.DecodeUint64(input)
	if errors.Is(err, hexutil.ErrMissingPrefix) {
		blckNum = cast.ToUint64(input)
	} else if err != nil {
		return err
	}

	if blckNum > math.MaxInt64 {
		return fmt.Errorf("block number larger than int64")
	}
	// #nosec G115 range checked
	*bn = BlockNumber(blckNum)

	return nil
}

// Int64 converts block number to primitive type
func (bn BlockNumber) Int64() int64 {
	if bn < 0 {
		return 0
	} else if bn == 0 {
		return 1
	}

	return int64(bn)
}

// TmHeight is a util function used for the Tendermint RPC client. It returns
// nil if the block number is "latest". Otherwise, it returns the pointer of the
// int64 value of the height.
func (bn BlockNumber) TmHeight() *int64 {
	if bn < 0 {
		return nil
	}

	height := bn.Int64()
	return &height
}

// BlockNumberOrHash represents a block number or a block hash.
type BlockNumberOrHash struct {
	BlockNumber *BlockNumber `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash `json:"blockHash,omitempty"`
}

func (bnh *BlockNumberOrHash) UnmarshalJSON(data []byte) error {
	type erased BlockNumberOrHash
	e := erased{}
	err := json.Unmarshal(data, &e)
	if err == nil {
		return bnh.checkUnmarshal(BlockNumberOrHash(e))
	}
	var input string
	err = json.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	err = bnh.decodeFromString(input)
	if err != nil {
		return err
	}

	return nil
}

func (bnh *BlockNumberOrHash) checkUnmarshal(e BlockNumberOrHash) error {
	if e.BlockNumber != nil && e.BlockHash != nil {
		return fmt.Errorf("cannot specify both BlockHash and BlockNumber, choose one or the other")
	}
	bnh.BlockNumber = e.BlockNumber
	bnh.BlockHash = e.BlockHash
	return nil
}

func (bnh *BlockNumberOrHash) decodeFromString(input string) error {
	switch input {
	case BlockParamEarliest:
		bn := EthEarliestBlockNumber
		bnh.BlockNumber = &bn
	case BlockParamLatest, BlockParamFinalized:
		bn := EthLatestBlockNumber
		bnh.BlockNumber = &bn
	case BlockParamPending:
		bn := EthPendingBlockNumber
		bnh.BlockNumber = &bn
	default:
		// check if the input is a block hash
		if len(input) == 66 {
			hash := common.Hash{}
			err := hash.UnmarshalText([]byte(input))
			if err != nil {
				return err
			}
			bnh.BlockHash = &hash
			break
		}
		// otherwise take the hex string has int64 value
		blockNumber, err := hexutil.DecodeUint64(input)
		if err != nil {
			return err
		}

		bnInt, err := ethermint.SafeInt64(blockNumber)
		if err != nil {
			return err
		}

		bn := BlockNumber(bnInt)
		bnh.BlockNumber = &bn
	}
	return nil
}

// https://github.com/ethereum/go-ethereum/blob/release/1.11/core/types/gen_header_json.go#L18
type Header struct {
	ParentHash common.Hash `json:"parentHash"       gencodec:"required"`
	UncleHash  common.Hash `json:"sha3Uncles"       gencodec:"required"`
	// update string avoid lost checksumed miner after MarshalText
	Coinbase    string              `json:"miner"`
	Root        common.Hash         `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash         `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash         `json:"receiptsRoot"     gencodec:"required"`
	Bloom       ethtypes.Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  *hexutil.Big        `json:"difficulty"       gencodec:"required"`
	Number      *hexutil.Big        `json:"number"           gencodec:"required"`
	GasLimit    hexutil.Uint64      `json:"gasLimit"         gencodec:"required"`
	GasUsed     hexutil.Uint64      `json:"gasUsed"          gencodec:"required"`
	Time        hexutil.Uint64      `json:"timestamp"        gencodec:"required"`
	Extra       hexutil.Bytes       `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash         `json:"mixHash"`
	Nonce       ethtypes.BlockNonce `json:"nonce"`
	BaseFee     *hexutil.Big        `json:"baseFeePerGas"                        rlp:"optional"`
	// overwrite rlpHash
	Hash common.Hash `json:"hash"`
}
