package e2etests

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestMigrateERC20CustodyFunds tests the migration of ERC20 custody funds
func TestMigrateERC20CustodyFunds(r *runner.E2ERunner, _ []string) {
	// get erc20 balance on ERC20 custody contract
	balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	require.NoError(r, err)

	// get EVM chain ID
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	newAddr := sample.EthAddress()

	// send MigrateERC20CustodyFunds command
	msg := crosschaintypes.NewMsgMigrateERC20CustodyFunds(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainID.Int64(),
		newAddr.Hex(),
		r.ERC20Addr.Hex(),
		sdkmath.NewUintFromBigInt(balance),
	)
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	event, ok := txserver.EventOfType[*crosschaintypes.EventERC20CustodyFundsMigration](res.Events)
	require.True(r, ok, "no EventERC20CustodyFundsMigration in %s", res.TxHash)

	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: event.CctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "migration")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(event.CctxIndex)

	// check ERC20 balance on new address
	newAddrBalance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, newAddr)
	require.NoError(r, err)
	require.Equal(r, balance, newAddrBalance)

	// artificially set the ERC20 Custody address to the new address to prevent accounting check from failing
	r.ERC20CustodyAddr = newAddr
}
