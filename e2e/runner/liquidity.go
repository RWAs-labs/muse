package runner

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/utils"
)

// AddLiquidityETH adds liquidity token to the uniswap pool MUSE/ETH
func (r *E2ERunner) AddLiquidityETH(amountMUSE, amountETH *big.Int) {
	r.ApproveETHMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.ETHMRC20Addr, amountMUSE, amountETH)
}

// AddLiquidityERC20 adds liquidity token to the uniswap pool MUSE/ERC20
func (r *E2ERunner) AddLiquidityERC20(amountMUSE, amountERC20 *big.Int) {
	r.ApproveERC20MRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.ERC20MRC20Addr, amountMUSE, amountERC20)
}

// AddLiquidityBTC adds liquidity token to the uniswap pool MUSE/BTC
func (r *E2ERunner) AddLiquidityBTC(amountMUSE, amountBTC *big.Int) {
	r.ApproveBTCMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.BTCMRC20Addr, amountMUSE, amountBTC)
}

// AddLiquiditySOL adds liquidity token to the uniswap pool MUSE/SOL
func (r *E2ERunner) AddLiquiditySOL(amountMUSE, amountSOL *big.Int) {
	r.ApproveSOLMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.SOLMRC20Addr, amountMUSE, amountSOL)
}

// AddLiquiditySPL adds liquidity token to the uniswap pool MUSE/SPL
func (r *E2ERunner) AddLiquiditySPL(amountMUSE, amountSPL *big.Int) {
	r.ApproveSPLMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.SPLMRC20Addr, amountMUSE, amountSPL)
}

// AddLiquiditySUI adds liquidity token to the uniswap pool MUSE/SUI
func (r *E2ERunner) AddLiquiditySUI(amountMUSE, amountSUI *big.Int) {
	r.ApproveSUIMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.SUIMRC20Addr, amountMUSE, amountSUI)
}

// AddLiquiditySuiFungibleToken adds liquidity token to the uniswap pool MUSE/SuiFungibleToken
func (r *E2ERunner) AddLiquiditySuiFungibleToken(amountMUSE, amountToken *big.Int) {
	r.ApproveFungibleTokenMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.SuiTokenMRC20Addr, amountMUSE, amountToken)
}

// AddLiquidityTON adds liquidity token to the uniswap pool MUSE/TON
func (r *E2ERunner) AddLiquidityTON(amountMUSE, amountTON *big.Int) {
	r.ApproveTONMRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.TONMRC20Addr, amountMUSE, amountTON)
}

// addLiquidity adds liquidity token to the uniswap pool MUSE/token
// we use the provided amount of MUSE and token to add liquidity as wanted amount
// 0 is used for the minimum amount of MUSE and token
func (r *E2ERunner) addLiquidity(tokenAddr ethcommon.Address, amountMUSE, amountToken *big.Int) {
	previousValue := r.MEVMAuth.Value
	r.MEVMAuth.Value = amountMUSE
	defer func() {
		r.MEVMAuth.Value = previousValue
	}()

	r.Logger.Info("Adding liquidity to MUSE/token pool")
	tx, err := r.UniswapV2Router.AddLiquidityETH(
		r.MEVMAuth,
		tokenAddr,
		amountToken,
		big.NewInt(0),
		big.NewInt(0),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == types.ReceiptStatusFailed {
		r.Logger.Error("Add liquidity failed for MUSE/token")
	}

	// get the pair address
	pairAddress, err := r.UniswapV2Factory.GetPair(&bind.CallOpts{}, r.WMuseAddr, tokenAddr)
	require.NoError(r, err)

	r.Logger.Info("MUSE/token pair address: %s", pairAddress.Hex())
}
