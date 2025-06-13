package observer

import (
	"context"
	stderrors "errors"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/compliance"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	contracts "github.com/RWAs-labs/muse/pkg/contracts/solana"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

var (
	gasOutboundParsers = []func(solana.CompiledInstruction) (contracts.OutboundInstruction, error){
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionWithdraw(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecute(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecuteRevert(inst)
		},
	}

	splOutboundParsers = []func(solana.CompiledInstruction) (contracts.OutboundInstruction, error){
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionWithdrawSPL(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecuteSPL(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecuteSPLRevert(inst)
		},
	}
)

// ProcessOutboundTrackers processes Solana outbound trackers
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.MusecoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "GetAllOutboundTrackerByChain error")
	}

	// prepare logger fields
	logger := ob.Logger().Outbound.With().
		Str("method", "ProcessOutboundTrackers").
		Int64("chain", chainID).
		Logger()

	// process outbound trackers
	for _, tracker := range trackers {
		// go to next tracker if this one already has a finalized tx
		nonce := tracker.Nonce
		if ob.IsTxFinalized(tracker.Nonce) {
			continue
		}

		// get original cctx parameters
		cctx, err := ob.MusecoreClient().GetCctxByNonce(ctx, chainID, tracker.Nonce)
		if err != nil {
			// take a rest if musecore RPC breaks
			return errors.Wrapf(err, "GetCctxByNonce error for chain %d nonce %d", chainID, tracker.Nonce)
		}
		coinType := cctx.InboundParams.CoinType

		// check each txHash and save its txResult if it's finalized and legit
		txCount := 0
		var txResult *rpc.GetTransactionResult
		for _, txHash := range tracker.HashList {
			if result, ok := ob.CheckFinalizedTx(ctx, txHash.TxHash, nonce, coinType); ok {
				txCount++
				txResult = result
				logger.Info().Msgf("confirmed outbound %s for chain %d nonce %d", txHash.TxHash, chainID, nonce)
				if txCount > 1 {
					logger.Error().
						Msgf("checkFinalizedTx passed, txCount %d chain %d nonce %d txResult %v", txCount, chainID, nonce, txResult)
				}
			}
		}
		// should be only one finalized txHash for each nonce
		if txCount == 1 {
			ob.SetTxResult(nonce, txResult)
		} else if txCount > 1 {
			// should not happen. We can't tell which txHash is true. It might happen (e.g. bug, glitchy/hacked endpoint)
			ob.Logger().Outbound.Error().Msgf("finalized multiple (%d) outbound for chain %d nonce %d", txCount, chainID, nonce)
		}
	}

	return nil
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *crosschaintypes.CrossChainTx) (bool, error) {
	// get outbound params
	params := cctx.GetCurrentOutboundParam()
	nonce := params.TssNonce
	coinType := cctx.InboundParams.CoinType

	// early return if outbound is not finalized yet
	txResult := ob.GetTxResult(nonce)
	if txResult == nil {
		return true, nil
	}

	// extract tx signature from tx result
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		// should not happen
		return false, errors.Wrapf(err, "GetTransaction error for nonce %d", nonce)
	}
	txSig := tx.Signatures[0]

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType)
	if err != nil {
		// should never happen as it was already successfully parsed in CheckFinalizedTx
		return false, errors.Wrapf(err, "ParseGatewayInstruction error for sig %s", txSig)
	}

	// the amount and status of the outbound
	outboundAmount := new(big.Int).SetUint64(inst.TokenAmount())

	// status was already verified as successful in CheckFinalizedTx
	outboundStatus := chains.ReceiveStatus_success
	if inst.InstructionDiscriminator() == contracts.DiscriminatorIncrementNonce {
		outboundStatus = chains.ReceiveStatus_failed
	}

	// cancelled transaction means the outbound is failed
	// - set amount to CCTX's amount to bypass amount check in musecore
	// - set status to failed to revert the CCTX in musecore
	if compliance.IsCCTXRestricted(cctx) {
		outboundAmount = cctx.GetCurrentOutboundParam().Amount.BigInt()
		outboundStatus = chains.ReceiveStatus_failed
	}

	// post vote to musecore
	ob.PostVoteOutbound(ctx, cctx.Index, txSig.String(), txResult, outboundAmount, outboundStatus, nonce, coinType)
	return false, nil
}

// PostVoteOutbound posts vote to musecore for the finalized outbound
func (ob *Observer) PostVoteOutbound(
	ctx context.Context,
	cctxIndex string,
	outboundHash string,
	txResult *rpc.GetTransactionResult,
	valueReceived *big.Int,
	status chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
) {
	// create outbound vote message
	msg := ob.CreateMsgVoteOutbound(cctxIndex, outboundHash, txResult, valueReceived, status, nonce, coinType)

	// prepare logger fields
	logFields := map[string]any{
		"chain": ob.Chain().ChainId,
		"nonce": nonce,
		"tx":    outboundHash,
	}

	// so we set retryGasLimit to 0 because the solana gateway withdrawal will always succeed
	// and the vote msg won't trigger MEVM interaction
	const (
		gasLimit = musecore.PostVoteOutboundGasLimit
	)

	retryGasLimit := musecore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = musecore.PostVoteOutboundRevertGasLimit
	}

	// post vote to musecore
	museTxHash, ballot, err := ob.MusecoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		ob.Logger().Outbound.Error().Err(err).Fields(logFields).Msg("PostVoteOutbound: error posting outbound vote")
		return
	}

	// print vote tx hash and ballot
	if museTxHash != "" {
		logFields["vote"] = museTxHash
		logFields["ballot"] = ballot
		ob.Logger().Outbound.Info().Fields(logFields).Msg("PostVoteOutbound: posted outbound vote successfully")
	}
}

// CreateMsgVoteOutbound creates a vote outbound message for Solana chain
func (ob *Observer) CreateMsgVoteOutbound(
	cctxIndex string,
	outboundHash string,
	txResult *rpc.GetTransactionResult,
	valueReceived *big.Int,
	status chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
) *crosschaintypes.MsgVoteOutbound {
	const (
		// Solana implements a different gas fee model than Ethereum, below values are not used.
		// Solana tx fee is based on both static fee and dynamic fee (priority fee), setting
		// zero values to by pass incorrectly funded gas stability pool.
		outboundGasUsed  = 0
		outboundGasPrice = 0
		outboundGasLimit = 0
	)

	creator := ob.MusecoreClient().GetKeys().GetOperatorAddress()

	return crosschaintypes.NewMsgVoteOutbound(
		creator.String(),
		cctxIndex,
		outboundHash,
		txResult.Slot, // instead of using block, Solana explorer uses slot for indexing
		outboundGasUsed,
		math.NewInt(outboundGasPrice),
		outboundGasLimit,
		math.NewUintFromBigInt(valueReceived),
		status,
		ob.Chain().ChainId,
		nonce,
		coinType,
		crosschaintypes.ConfirmationMode_SAFE,
	)
}

// CheckFinalizedTx checks if a txHash is finalized for given nonce and coinType
// returns (tx result, true) if finalized or (nil, false) otherwise
func (ob *Observer) CheckFinalizedTx(
	ctx context.Context,
	txHash string,
	nonce uint64,
	coinType coin.CoinType,
) (*rpc.GetTransactionResult, bool) {
	// prepare logger fields
	chainID := ob.Chain().ChainId
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "CheckFinalizedTx").
		Int64(logs.FieldChain, chainID).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).Logger()

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		logger.Error().Err(err).Msg("SignatureFromBase58 error")
		return nil, false
	}

	// query transaction using "finalized" commitment to avoid re-org
	txResult, err := ob.solClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		logger.Error().Err(err).Msg("GetTransaction error")
		return nil, false
	}

	// the tx must be successful in order to effectively increment the nonce
	if txResult.Meta.Err != nil {
		logger.Error().Any("Err", txResult.Meta.Err).Msg("tx is not successful")
		return nil, false
	}

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType)
	if err != nil {
		logger.Error().Err(err).Msg("ParseGatewayInstruction error")
		return nil, false
	}

	txNonce := inst.GatewayNonce()

	// recover ECDSA signer from instruction
	signerECDSA, err := inst.Signer()
	if err != nil {
		logger.Error().Err(err).Msg("cannot get instruction signer")
		return nil, false
	}

	// check tx authorization
	if signerECDSA != ob.TSS().PubKey().AddressEVM() {
		logger.Error().
			Msgf("tx signer %s is not matching current TSS address %s", signerECDSA, ob.TSS().PubKey().AddressEVM())
		return nil, false
	}

	// check tx nonce
	if txNonce != nonce {
		logger.Error().Msgf("tx nonce %d is not matching tracker nonce", txNonce)
		return nil, false
	}

	return txResult, true
}

// parseInstructionWith attempts to parse an instruction using a list of parsers
func parseInstructionWith(
	instruction solana.CompiledInstruction,
	parsers []func(solana.CompiledInstruction) (contracts.OutboundInstruction, error),
) (contracts.OutboundInstruction, error) {
	errs := make([]error, 0, len(parsers))
	for _, parser := range parsers {
		inst, err := parser(instruction)
		if err == nil {
			return inst, nil
		}
		errs = append(errs, err)
	}
	return nil, errors.Wrap(stderrors.Join(errs...), "failed to parse instruction")
}

// ParseGatewayInstruction parses the outbound instruction from tx result
func ParseGatewayInstruction(
	txResult *rpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	coinType coin.CoinType,
) (contracts.OutboundInstruction, error) {
	// unmarshal transaction
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling transaction")
	}

	// validate instruction count
	// if there are 2 instructions, first one can only be optional compute budget instruction
	instructionCount := len(tx.Message.Instructions)
	if instructionCount < 1 || instructionCount > 2 {
		return nil, fmt.Errorf("unexpected number of instructions: %d", instructionCount)
	}

	// get gateway instruction
	instruction := tx.Message.Instructions[instructionCount-1]

	// if there are two instructions, validate the first one program is compute budget
	if instructionCount == 2 {
		budgetProgramID, err := tx.Message.Program(tx.Message.Instructions[0].ProgramIDIndex)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve program ID")
		}

		if !budgetProgramID.Equals(solana.ComputeBudget) {
			return nil, fmt.Errorf(
				"programID %s is not matching compute budget id %s",
				budgetProgramID,
				solana.ComputeBudget,
			)
		}
	}

	// validate gateway instruction program
	programID, err := tx.Message.Program(instruction.ProgramIDIndex)
	if err != nil {
		return nil, errors.Wrap(err, "error getting program ID")
	}

	// the instruction should be an invocation of the gateway program
	if !programID.Equals(gatewayID) {
		return nil, fmt.Errorf("programID %s is not matching gatewayID %s", programID, gatewayID)
	}

	// first check if it was simple nonce increment instruction, which indicates that outbound failed
	inst, err := contracts.ParseInstructionIncrementNonce(instruction)
	if err == nil {
		return inst, nil
	}

	// parse the outbound instruction
	switch coinType {
	case coin.CoinType_Gas:
		return parseInstructionWith(instruction, gasOutboundParsers)
	case coin.CoinType_ERC20:
		return parseInstructionWith(instruction, splOutboundParsers)
	case coin.CoinType_Cmd:
		return contracts.ParseInstructionWhitelist(instruction)
	case coin.CoinType_NoAssetCall:
		return contracts.ParseInstructionExecute(instruction)
	default:
		return nil, fmt.Errorf("unsupported outbound coin type %s", coinType)
	}
}
