package signer

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"

	contracts "github.com/RWAs-labs/muse/pkg/contracts/solana"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// prepareWhitelistTx prepares whitelist outbound
func (signer *Signer) prepareWhitelistTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
) (outboundGetter, error) {
	params := cctx.GetCurrentOutboundParam()
	relayedMsg := strings.Split(cctx.RelayedMessage, ":")
	if len(relayedMsg) != 2 {
		return nil, fmt.Errorf("TryProcessOutbound: invalid relayed msg")
	}

	pk, err := solana.PublicKeyFromBase58(relayedMsg[1])
	if err != nil {
		return nil, errors.Wrapf(err, "publicKeyFromBase58 %s error", relayedMsg[1])
	}

	seed := [][]byte{[]byte("whitelist"), pk.Bytes()}
	whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, signer.gatewayID)
	if err != nil {
		return nil, errors.Wrapf(err, "findProgramAddress error for seed %s", seed)
	}

	// sign gateway whitelist message by TSS
	msg, err := signer.createAndSignMsgWhitelist(ctx, params, height, pk, whitelistEntryPDA)
	if err != nil {
		return nil, errors.Wrap(err, "createAndSignMsgWhitelist error")
	}

	return func() (*Outbound, error) {
		// sign the whitelist transaction by relayer key
		inst, err := signer.createWhitelistInstruction(msg)
		if err != nil {
			return nil, errors.Wrap(err, "error creating whitelist instruction")
		}

		tx, err := signer.signTx(ctx, inst, 0)
		if err != nil {
			return nil, errors.Wrap(err, "error signing whitelist instruction")
		}
		return &Outbound{Tx: tx}, nil
	}, nil
}

// createAndSignMsgWhitelist creates and signs a whitelist message (for gateway whitelist_spl_mint instruction) with TSS.
func (signer *Signer) createAndSignMsgWhitelist(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
	whitelistCandidate solana.PublicKey,
	whitelistEntry solana.PublicKey,
) (*contracts.MsgWhitelist, error) {
	chain := signer.Chain()
	// #nosec G115 always positive
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce

	// prepare whitelist msg and compute hash
	msg := contracts.NewMsgWhitelist(whitelistCandidate, whitelistEntry, chainID, nonce)
	msgHash := msg.Hash()

	// sign the message with TSS to get an ECDSA signature.
	// the produced signature is in the [R || S || V] format where V is 0 or 1.
	signature, err := signer.TSS().Sign(ctx, msgHash[:], height, nonce, chain.ChainId)
	if err != nil {
		return nil, errors.Wrap(err, "key-sign failed")
	}

	// attach the signature and return
	return msg.SetSignature(signature), nil
}

// createWhitelistInstruction wraps the whitelist 'msg' into a Solana instruction.
func (signer *Signer) createWhitelistInstruction(msg *contracts.MsgWhitelist) (*solana.GenericInstruction, error) {
	// create whitelist_spl_mint instruction with program call data
	dataBytes, err := borsh.Serialize(contracts.WhitelistInstructionParams{
		Discriminator: contracts.DiscriminatorWhitelistSplMint,
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize whitelist_spl_mint instruction")
	}

	inst := &solana.GenericInstruction{
		ProgID:    signer.gatewayID,
		DataBytes: dataBytes,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer.relayerKey.PublicKey()).WRITE().SIGNER(),
			solana.Meta(signer.pda).WRITE(),
			solana.Meta(msg.WhitelistEntry()).WRITE(),
			solana.Meta(msg.WhitelistCandidate()),
			solana.Meta(solana.SystemProgramID),
		},
	}

	return inst, nil
}
