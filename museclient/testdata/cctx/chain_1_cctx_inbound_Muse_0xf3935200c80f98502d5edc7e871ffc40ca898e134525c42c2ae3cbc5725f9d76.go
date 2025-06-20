package cctx

import (
	sdkmath "cosmossdk.io/math"

	"github.com/RWAs-labs/muse/pkg/coin"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

var chain_1_cctx_inbound_Muse_0xf393520 = &crosschaintypes.CrossChainTx{
	Creator:        "muse1p0uwsq4naus5r4l7l744upy0k8ezzj84mn40nf",
	Index:          "0x477544c4b8c8be544b23328b21286125c89cd6bb5d1d6d388d91eea8ea1a6f1f",
	MuseFees:       sdkmath.NewUintFromString("0"),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "Remote omnichain contract call completed",
		LastUpdateTimestamp: 1708490549,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0x2f993766e8e1Ef9288B1F33F6aa244911A0A77a7",
		SenderChainId:          1,
		TxOrigin:               "0x2f993766e8e1Ef9288B1F33F6aa244911A0A77a7",
		CoinType:               coin.CoinType_Muse,
		Asset:                  "",
		Amount:                 sdkmath.NewUintFromString("20000000000000000000"),
		ObservedHash:           "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76",
		ObservedExternalHeight: 19273702,
		BallotIndex:            "0x611f8fea15d26d318c06b04e41bcc16ef212048eddbbf62641c1e3e42b376009",
		FinalizedMuseHeight:    1851403,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:        "0x2f993766e8e1ef9288b1f33f6aa244911a0a77a7",
			ReceiverChainId: 7000,
			CoinType:        coin.CoinType_Muse,
			Amount:          sdkmath.ZeroUint(),
			TssNonce:        0,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 100000,
			},
			GasPrice:               "",
			Hash:                   "0x947434364da7c74d7e896a389aa8cb3122faf24bbcba64b141cb5acd7838209c",
			BallotIndex:            "",
			ObservedExternalHeight: 1851403,
			GasUsed:                0,
			EffectiveGasPrice:      sdkmath.ZeroInt(),
			EffectiveGasLimit:      0,
			TssPubkey:              "musepub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
		},
	},
}
