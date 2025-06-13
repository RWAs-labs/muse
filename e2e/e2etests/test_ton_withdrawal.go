package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/museclient/chains/ton/liteapi"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
)

func TestTONWithdraw(r *runner.E2ERunner, args []string) {
	// ARRANGE
	require.Len(r, args, 1)

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given mEVM sender
	mevmSender := r.MEVMAuth.From

	// Given his MRC-20 balance
	senderMRC20BalanceBefore, err := r.TONMRC20.BalanceOf(&bind.CallOpts{}, mevmSender)
	require.NoError(r, err)
	r.Logger.Info("mEVM sender's MRC20 TON balance before withdraw: %d", senderMRC20BalanceBefore)

	// Given a receiver
	_, receiver, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	receiverBalanceBefore, err := r.Clients.TON.GetBalanceOf(r.Ctx, receiver.GetAddress(), true)
	require.NoError(r, err)

	r.Logger.Info("Recipient's TON balance before withdrawal: %s", toncontracts.FormatCoins(receiverBalanceBefore))
	r.Logger.Info("Receiver's TON address: %s", receiver.GetAddress().ToHuman(false, true))

	// Given amount to withdraw (and approved amount in TON MRC20 to cover the gas fee)
	amount := utils.ParseUint(r, args[0])
	r.Logger.Info("Amount to withdraw: %s", toncontracts.FormatCoins(amount))
	approvedAmount := amount.Add(toncontracts.Coins(1))

	// ACT
	cctx := r.WithdrawTONMRC20(receiver.GetAddress(), amount.BigInt(), approvedAmount.BigInt())

	// ASSERT
	r.Logger.Info(
		"Withdraw TON MRC20 transaction (with %s) sent: %+v",
		toncontracts.FormatCoins(amount),
		map[string]any{
			"mevm_sender":   mevmSender.Hex(),
			"ton_recipient": receiver.GetAddress().ToRaw(),
			"ton_amount":    toncontracts.FormatCoins(amount),
			"cctx_index":    cctx.Index,
			"ton_hash":      cctx.GetCurrentOutboundParam().Hash,
			"mevm_hash":     cctx.InboundParams.ObservedHash,
		},
	)

	// Make sure that recipient's TON balance has increased
	receiverBalanceAfter, err := r.Clients.TON.GetBalanceOf(r.Ctx, receiver.GetAddress(), true)
	require.NoError(r, err)

	r.Logger.Info("Recipient's balance after withdrawal: %s", toncontracts.FormatCoins(receiverBalanceAfter))

	// Make sure that sender's MRC20 balance has decreased
	senderMRC20BalanceAfter, err := r.TONMRC20.BalanceOf(&bind.CallOpts{}, mevmSender)
	require.NoError(r, err)
	r.Logger.Info("mEVM sender's MRC20 TON balance after withdraw: %d", senderMRC20BalanceAfter)
	r.Logger.Info(
		"mEVM sender's MRC20 TON balance diff: %d",
		big.NewInt(0).Sub(senderMRC20BalanceBefore, senderMRC20BalanceAfter),
	)

	// Make sure that TON withdrawal CCTX contain outgoing message with exact withdrawal amount
	lt, hash, err := liteapi.TransactionHashFromString(cctx.GetCurrentOutboundParam().Hash)
	require.NoError(r, err)

	txs, err := r.Clients.TON.GetTransactions(r.Ctx, 1, gw.AccountID(), lt, hash)
	require.NoError(r, err)
	require.Len(r, txs, 1)

	// TON coins that were withdrawn from GW to the recipient
	inMsgAmount := math.NewUint(
		uint64(txs[0].Msgs.OutMsgs.Values()[0].Value.Info.IntMsgInfo.Value.Grams),
	)

	// #nosec G115 always in range
	require.Equal(r, int(amount.Uint64()), int(inMsgAmount.Uint64()))
}
