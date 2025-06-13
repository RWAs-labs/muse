package e2etests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/erc20"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestWhitelistERC20 tests the whitelist ERC20 functionality
func TestWhitelistERC20(r *runner.E2ERunner, _ []string) {
	// Deploy a new ERC20 on the new EVM chain
	r.Logger.Info("Deploying new ERC20 contract")
	erc20Addr, txERC20, _, err := erc20.DeployERC20(r.EVMAuth, r.EVMClient, "NEWERC20", "NEWERC20", 6)
	require.NoError(r, err)

	// wait for the ERC20 to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txERC20, r.Logger, r.ReceiptTimeout)
	require.Equal(r, ethtypes.ReceiptStatusSuccessful, receipt.Status)

	// ERC20 test

	// whitelist erc20 mrc20
	r.Logger.Info("whitelisting ERC20 on new network")
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, crosschaintypes.NewMsgWhitelistERC20(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		erc20Addr.Hex(),
		chains.GoerliLocalnet.ChainId,
		"NEWERC20",
		"NEWERC20",
		6,
		100000,
		sdkmath.NewUintFromString("100000000000000000000000000"),
	))
	require.NoError(r, err)

	event, ok := txserver.EventOfType[*crosschaintypes.EventERC20Whitelist](res.Events)
	require.True(r, ok, "no EventERC20Whitelist in %s", res.TxHash)
	erc20mrc20Addr := event.Mrc20Address
	whitelistCCTXIndex := event.WhitelistCctxIndex

	err = r.MuseTxServer.InitializeLiquidityCaps(erc20mrc20Addr)
	require.NoError(r, err)

	// ensure CCTX created
	resCCTX, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: whitelistCCTXIndex})
	require.NoError(r, err)

	cctx := resCCTX.CrossChainTx
	r.Logger.CCTX(*cctx, "whitelist_cctx")

	// wait for the whitelist cctx to be mined
	r.WaitForMinedCCTXFromIndex(whitelistCCTXIndex)

	// save old ERC20 attribute to set it back after the test
	oldERC20Addr := r.ERC20Addr
	oldERC20 := r.ERC20
	oldERC20MRC20Addr := r.ERC20MRC20Addr
	oldERC20MRC20 := r.ERC20MRC20
	defer func() {
		r.ERC20Addr = oldERC20Addr
		r.ERC20 = oldERC20
		r.ERC20MRC20Addr = oldERC20MRC20Addr
		r.ERC20MRC20 = oldERC20MRC20
	}()

	// set erc20 and mrc20 in runner
	require.True(r, ethcommon.IsHexAddress(erc20mrc20Addr), "invalid contract address: %s", erc20mrc20Addr)
	erc20mrc20AddrHex := ethcommon.HexToAddress(erc20mrc20Addr)
	erc20MRC20, err := mrc20.NewMRC20(erc20mrc20AddrHex, r.MEVMClient)
	require.NoError(r, err)
	r.ERC20MRC20Addr = erc20mrc20AddrHex
	r.ERC20MRC20 = erc20MRC20

	erc20ERC20, err := erc20.NewERC20(erc20Addr, r.EVMClient)
	require.NoError(r, err)
	r.ERC20Addr = erc20Addr
	r.ERC20 = erc20ERC20

	// get balance
	balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.Account.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("ERC20 balance: %s", balance.String())

	// run deposit and withdraw ERC20 test
	txHash := r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), balance, []byte{})
	r.WaitForMinedCCTX(txHash)

	// approve 1 unit of the gas token to cover the gas fee
	tx, err := r.ETHMRC20.Approve(r.MEVMAuth, r.ERC20MRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Info("eth mrc20 approve receipt: status %d", receipt.Status)

	tx = r.LegacyWithdrawERC20(balance)
	r.WaitForMinedCCTX(tx.Hash())
}
