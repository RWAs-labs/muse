package cctx

import (
	sdkmath "cosmossdk.io/math"

	"github.com/RWAs-labs/muse/pkg/coin"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

var chain_1_cctx_inbound_ERC20_0x4ea69a0 = &crosschaintypes.CrossChainTx{
	Creator:        "muse1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
	Index:          "0xd326700a1931f28853f44f8462f72588f94b1f248214d59a23c3e1b141ff5ede",
	MuseFees:       sdkmath.NewUintFromString("0"),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "Remote omnichain contract call completed",
		LastUpdateTimestamp: 1709052990,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0x56BF8D4a6E7b59D2C0E40Cba2409a4a30ab4FbE2",
		SenderChainId:          1,
		TxOrigin:               "0x56BF8D4a6E7b59D2C0E40Cba2409a4a30ab4FbE2",
		CoinType:               coin.CoinType_ERC20,
		Asset:                  "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Amount:                 sdkmath.NewUintFromString("9992000000"),
		ObservedHash:           "0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da",
		ObservedExternalHeight: 19320188,
		BallotIndex:            "0x97101937e3927e124dffcaed1349af2599a8420ff34315288e96eac7f0033048",
		FinalizedMuseHeight:    1944675,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:        "0x56bf8d4a6e7b59d2c0e40cba2409a4a30ab4fbe2",
			ReceiverChainId: 7000,
			CoinType:        coin.CoinType_ERC20,
			Amount:          sdkmath.NewUintFromString("0"),
			TssNonce:        0,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 1500000,
			},
			GasPrice:               "",
			Hash:                   "0xf63eaa3e01af477673aa9e86fb634df15d30a00734dab7450cb0fc28dbc9d11b",
			BallotIndex:            "",
			ObservedExternalHeight: 1944675,
			GasUsed:                0,
			EffectiveGasPrice:      sdkmath.NewInt(0),
			EffectiveGasLimit:      0,
			TssPubkey:              "musepub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
		},
	},
}
