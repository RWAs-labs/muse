package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/museclient/chains/ton/liteapi"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func TestTONDepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	getBalance := func(addr ethcommon.Address) uint64 {
		b, err := r.TONMRC20.BalanceOf(&bind.CallOpts{}, addr)
		require.NoError(r, err)

		return b.Uint64()
	}

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Given restricted receiver ...
	recipient := ethcommon.HexToAddress(sample.RestrictedEVMAddressTest)

	// ... and its balance before deposit
	oldBalance := getBalance(recipient)

	// ACT
	tx, err := r.TONDepositRaw(gw, sender, amount, recipient)
	require.NoError(r, err)

	tonHash := liteapi.TransactionToHashString(tx)

	r.WaitForBlocks(5)

	// ASSERT
	// No cctx was created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tonHash, r.CctxClient)

	// Receiver's balance IS NOT changed
	newBalance := getBalance(recipient)
	require.Equal(r, oldBalance, newBalance)
}
