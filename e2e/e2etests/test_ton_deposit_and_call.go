package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/example"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
)

func TestTONDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx depositAndCall fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDepositAndCall)
	require.NoError(r, err)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Given sample mEVM contract deployed by userTON account
	contractAddr, _, contract, err := testcontract.DeployExample(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err, "unable to deploy example contract")
	r.Logger.Info("Example mevm contract deployed at: %s", contractAddr.String())

	// Given call data
	callData := []byte("hello from TON!")

	// ACT
	_, err = r.TONDepositAndCall(gw, sender, amount, contractAddr, callData)

	// ASSERT
	require.NoError(r, err)

	expectedDeposit := amount.Sub(depositFee)

	// check if example contract has been called, bar value should be set to amount
	utils.MustHaveCalledExampleContract(r, contract, expectedDeposit.BigInt(), []byte(sender.GetAddress().ToRaw()))

	// Check receiver's balance
	balance, err := r.TONMRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)

	r.Logger.Info("Contract's mEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
