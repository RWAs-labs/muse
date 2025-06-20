package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// ValidateInbound is the only entry-point to create new CCTX (eg. when observers voting is done or new inbound event is detected).
// It creates new CCTX object and calls InitiateOutbound method.
func (k Keeper) ValidateInbound(
	ctx sdk.Context,
	msg *types.MsgVoteInbound,
	shouldPayGas bool,
) (*types.CrossChainTx, error) {
	tss, tssFound := k.museObserverKeeper.GetTSS(ctx)
	if !tssFound {
		return nil, types.ErrCannotFindTSSKeys
	}
	if err := k.CheckIfTSSMigrationTransfer(ctx, msg); err != nil {
		return nil, errors.Wrap(err, "tss migration transfer check failed")
	}

	// Do not process if inbound is disabled
	if !k.museObserverKeeper.IsInboundEnabled(ctx) {
		return nil, observertypes.ErrInboundDisabled
	}

	// create a new CCTX from the inbound message. The status of the new CCTX is set to PendingInbound.
	cctx, err := types.NewCCTX(ctx, *msg, tss.TssPubkey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new CCTX")
	}

	// Initiate outbound, the process function manages the state commit and cctx status change.
	// If the process fails, the changes to the evm state are rolled back.
	_, err = k.InitiateOutbound(ctx, InitiateOutboundConfig{
		CCTX:         &cctx,
		ShouldPayGas: shouldPayGas,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initiate outbound")
	}

	inCctxIndex, ok := ctx.Value(InCCTXIndexKey).(string)
	if ok {
		cctx.InboundParams.ObservedHash = inCctxIndex
	}
	k.SaveCCTXUpdate(ctx, cctx, tss.TssPubkey)

	return &cctx, nil
}

// CheckIfTSSMigrationTransfer checks if the sender is a TSS address and returns an error if it is.
// If the sender is an older TSS address, this means that it is a migration transfer, and we do not need to treat this as a deposit and process the CCTX
func (k Keeper) CheckIfTSSMigrationTransfer(ctx sdk.Context, msg *types.MsgVoteInbound) error {
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	historicalTSSList := k.museObserverKeeper.GetAllTSS(ctx)
	chain, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId)
	if !found {
		return observertypes.ErrSupportedChains.Wrapf("chain not found for chainID %d", msg.SenderChainId)
	}

	// the check is only necessary if the inbound is validated from observers from a connected chain
	if chain.CctxGateway != chains.CCTXGateway_observers {
		return nil
	}

	switch {
	case chains.IsEVMChain(chain.ChainId, additionalChains):
		for _, tss := range historicalTSSList {
			ethTssAddress, err := crypto.GetTSSAddrEVM(tss.TssPubkey)
			if err != nil {
				continue
			}
			if ethTssAddress.Hex() == msg.Sender {
				return types.ErrMigrationFromOldTss
			}
		}
	case chains.IsBitcoinChain(chain.ChainId, additionalChains):
		bitcoinParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
		if err != nil {
			return err
		}
		for _, tss := range historicalTSSList {
			btcTssAddress, err := crypto.GetTSSAddrBTC(tss.TssPubkey, bitcoinParams)
			if err != nil {
				continue
			}
			if btcTssAddress == msg.Sender {
				return types.ErrMigrationFromOldTss
			}
		}
	}

	return nil
}
