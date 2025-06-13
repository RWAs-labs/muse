package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// TestUpdateMRC20Name tests updating name and symbol of a MRC20
func TestUpdateMRC20Name(r *runner.E2ERunner, _ []string) {
	msg := fungibletypes.NewMsgUpdateMRC20Name(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ETHMRC20Addr.Hex(),
		"New ETH",
		"ETH.NEW",
	)
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update eth mrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the MRC20
	newName, err := r.ETHMRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "New ETH", newName)

	newSymbol, err := r.ETHMRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "ETH.NEW", newSymbol)

	qRes, err := r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHMRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.EqualValues(r, "New ETH", qRes.ForeignCoins.Name)
	require.EqualValues(r, "ETH.NEW", qRes.ForeignCoins.Symbol)

	// try another mrc20
	msg = fungibletypes.NewMsgUpdateMRC20Name(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ERC20MRC20Addr.Hex(),
		"New USDT",
		"USDT.NEW",
	)
	res, err = r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update erc20 mrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the MRC20
	newName, err = r.ERC20MRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "New USDT", newName)

	newSymbol, err = r.ERC20MRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "USDT.NEW", newSymbol)

	qRes, err = r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ERC20MRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.EqualValues(r, "New USDT", qRes.ForeignCoins.Name)
	require.EqualValues(r, "USDT.NEW", qRes.ForeignCoins.Symbol)
}
