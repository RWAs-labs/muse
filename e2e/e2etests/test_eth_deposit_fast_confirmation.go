package e2etests

import (
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	mathpkg "github.com/RWAs-labs/muse/pkg/math"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// TestETHDepositFastConfirmation tests the fast confirmation of ETH deposits
func TestETHDepositFastConfirmation(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// ARRANGE
	// query chainID
	chainIDBig, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	chainID := chainIDBig.Int64()

	// enable inbound fast confirmation by updating the chain params
	reqQuery := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}
	resOldChainParams, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, reqQuery)
	require.NoError(r, err)

	chainParams := *resOldChainParams.ChainParams
	chainParams.ConfirmationParams = &observertypes.ConfirmationParams{
		SafeInboundCount:  10, // approx 10 seconds, much longer than Fast confirmation time (1 second)
		FastInboundCount:  1,
		SafeOutboundCount: 1,
		FastOutboundCount: 1,
	}
	err = r.MuseTxServer.UpdateChainParams(&chainParams)
	require.NoError(r, err, "failed to enable inbound fast confirmation")

	// it takes 1 Muse block time for museclient to pick up the new chain params
	// wait for 2 blocks to ensure the new chain params are effective
	utils.WaitForMuseBlocks(r.Ctx, r, r.MEVMClient, 2, 20*time.Second)
	r.Logger.Info("enabled inbound fast confirmation")

	// query current ETH MRC20 supply
	supply, err := r.ETHMRC20.TotalSupply(&bind.CallOpts{})
	supplyUint := sdkmath.NewUintFromBigInt(supply)
	require.NoError(r, err)

	// set MRC20 liquidity cap to 150% of the current supply
	// note: the percentage should not be too small as it may block other tests
	liquidityCap, _ := mathpkg.IncreaseUintByPercent(supplyUint, 50)
	require.True(r, liquidityCap.GT(sdkmath.ZeroUint()))
	res, err := r.MuseTxServer.SetMRC20LiquidityCap(r.ETHMRC20Addr, liquidityCap)
	require.NoError(r, err)
	r.Logger.Info("set liquidity cap to %s tx hash: %s", liquidityCap.String(), res.TxHash)

	// ACT-1
	// deposit with exactly fast amount cap, should be fast confirmed
	fastAmountCap := chains.CalcInboundFastConfirmationAmountCap(chainID, liquidityCap)
	tx := r.ETHDeposit(
		r.EVMAddress(),
		fastAmountCap.BigInt(),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		true,
	)
	r.Logger.Info("deposited exactly fast amount %d cap tx hash: %s", fastAmountCap, tx.Hash().Hex())

	// ASSERT-1
	// wait for the cctx to be FAST confirmed
	timeStart := time.Now()
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(r, crosschaintypes.ConfirmationMode_FAST, cctx.InboundParams.ConfirmationMode)
	fastConfirmTime := time.Since(timeStart)

	r.Logger.Info("FAST confirmed deposit succeeded in %f seconds", fastConfirmTime.Seconds())

	// ACT-2
	// deposit with amount more than fast amount cap
	amountMoreThanCap := big.NewInt(0).Add(fastAmountCap.BigInt(), big.NewInt(1))
	tx = r.ETHDeposit(
		r.EVMAddress(),
		amountMoreThanCap,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		true,
	)
	r.Logger.Info("deposited more than fast amount cap %d tx hash: %s", amountMoreThanCap, tx.Hash().Hex())

	// ASSERT-2
	// wait for the cctx to be SAFE confirmed
	timeStart = time.Now()
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(r, crosschaintypes.ConfirmationMode_SAFE, cctx.InboundParams.ConfirmationMode)
	safeConfirmTime := time.Since(timeStart)

	r.Logger.Info("SAFE confirmed deposit succeeded in %f seconds", safeConfirmTime.Seconds())

	// ensure FAST confirmation is faster than SAFE confirmation
	// using 3 seconds is good enough to check the difference on local goerli network
	timeSaved := safeConfirmTime - fastConfirmTime
	r.Logger.Info("FAST confirmation saved %f seconds", timeSaved.Seconds())
	require.True(r, timeSaved > 3*time.Second)

	// TEARDOWN
	// restore old chain params
	err = r.MuseTxServer.UpdateChainParams(resOldChainParams.ChainParams)
	require.NoError(r, err, "failed to restore chain params")

	// remove the liquidity cap
	_, err = r.MuseTxServer.RemoveMRC20LiquidityCap(r.ETHMRC20Addr)
	require.NoError(r, err)
}
