package chains

import (
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/cmd/musetool/context"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	museclientObserver "github.com/RWAs-labs/muse/museclient/chains/bitcoin/observer"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/memo"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func BitcoinBallotIdentifier(
	ctx *context.Context,
	btcClient *client.Client,
	params *chaincfg.Params,
	tss string,
	txHash string,
	senderChainID int64,
	musecoreChainID int64,
	confirmationCount uint64,
) (cctxIdentifier string, isConfirmed bool, err error) {
	var (
		goCtx = ctx.GetContext()
	)

	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return
	}
	tx, err := btcClient.GetRawTransactionVerbose(goCtx, hash)
	if err != nil {
		return
	}

	if tx.Confirmations >= confirmationCount {
		isConfirmed = true
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return
	}

	blockVb, err := btcClient.GetBlockVerbose(goCtx, blockHash)
	if err != nil {
		return
	}

	event, err := museclientObserver.GetBtcEventWithWitness(
		goCtx,
		btcClient,
		*tx,
		tss,
		uint64(blockVb.Height), // #nosec G115 always positive
		zerolog.New(zerolog.Nop()),
		params,
		common.CalcDepositorFee,
	)
	if err != nil {
		return
	}
	if event == nil {
		err = fmt.Errorf("no event built for btc sent to TSS")
		return
	}

	cctxIdentifier, err = identifierFromBtcEvent(event, senderChainID, musecoreChainID)
	return
}

func identifierFromBtcEvent(event *museclientObserver.BTCInboundEvent,
	senderChainID int64,
	musecoreChainID int64) (cctxIdentifier string, err error) {
	// decode event memo bytes
	err = event.DecodeMemoBytes(senderChainID)
	if err != nil {
		return
	}

	// convert the amount to integer (satoshis)
	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		return
	}
	amountInt := big.NewInt(amountSats)

	var msg *crosschaintypes.MsgVoteInbound
	switch event.MemoStd {
	case nil:
		{
			msg = voteFromLegacyMemo(event, amountInt, senderChainID, musecoreChainID)
		}
	default:
		{
			msg = voteFromStdMemo(event, amountInt, senderChainID, musecoreChainID)
		}
	}
	if msg == nil {
		return
	}

	cctxIdentifier = msg.Digest()
	return
}

// NewInboundVoteFromLegacyMemo creates a MsgVoteInbound message for inbound that uses legacy memo
func voteFromLegacyMemo(
	event *museclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	musecoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	message := hex.EncodeToString(event.MemoBytes)

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.FromAddress,
		senderChainID,
		event.FromAddress,
		event.ToAddress,
		musecoreChainID,
		cosmosmath.NewUintFromBigInt(amountSats),
		message,
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not relevant for v1
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithCrossChainCall(len(event.MemoBytes) > 0),
	)
}

func voteFromStdMemo(
	event *museclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	musecoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	// musecore will create a revert outbound that points to the custom revert address.
	revertOptions := crosschaintypes.RevertOptions{
		RevertAddress: event.MemoStd.RevertOptions.RevertAddress,
	}

	// check if the memo is a cross-chain call, or simple token deposit
	isCrosschainCall := event.MemoStd.OpCode == memo.OpCodeCall || event.MemoStd.OpCode == memo.OpCodeDepositAndCall

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.FromAddress,
		senderChainID,
		event.FromAddress,
		event.ToAddress,
		musecoreChainID,
		cosmosmath.NewUintFromBigInt(amountSats),
		hex.EncodeToString(event.MemoStd.Payload),
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not relevant for v1
		event.Status,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithRevertOptions(revertOptions),
		crosschaintypes.WithCrossChainCall(isCrosschainCall),
	)
}
