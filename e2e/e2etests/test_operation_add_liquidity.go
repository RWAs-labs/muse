package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
)

func parseArgsForAddLiquidity(r *runner.E2ERunner, args []string) (*big.Int, *big.Int) {
	require.Len(r, args, 2)

	liqMUSE := utils.ParseBigInt(r, args[0])
	liqToken := utils.ParseBigInt(r, args[1])

	return liqMUSE, liqToken
}

// TestOperationAddLiquidityETH is an operational test to add liquidity in the MUSE/ETH pool (evm gas token)
func TestOperationAddLiquidityETH(r *runner.E2ERunner, args []string) {
	liqMUSE, liqETH := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityETH(liqMUSE, liqETH)
}

// TestOperationAddLiquidityERC20 is an operational test to add liquidity in the MUSE/ERC20 pool
func TestOperationAddLiquidityERC20(r *runner.E2ERunner, args []string) {
	liqMUSE, liqERC20 := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityERC20(liqMUSE, liqERC20)
}

// TestOperationAddLiquidityBTC is an operational test to add liquidity in the MUSE/BTC pool
func TestOperationAddLiquidityBTC(r *runner.E2ERunner, args []string) {
	liqMUSE, liqBTC := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityBTC(liqMUSE, liqBTC)
}

// TestOperationAddLiquiditySOL is an operational test to add liquidity in the MUSE/SOL pool
func TestOperationAddLiquiditySOL(r *runner.E2ERunner, args []string) {
	liqMUSE, liqSOL := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySOL(liqMUSE, liqSOL)
}

// TestOperationAddLiquiditySPL is an operational test to add liquidity in the MUSE/SPL pool
func TestOperationAddLiquiditySPL(r *runner.E2ERunner, args []string) {
	liqMUSE, liqSPL := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySPL(liqMUSE, liqSPL)
}

// TestOperationAddLiquiditySUI is an operational test to add liquidity in the MUSE/SUI pool
func TestOperationAddLiquiditySUI(r *runner.E2ERunner, args []string) {
	liqMUSE, liqSUI := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySUI(liqMUSE, liqSUI)
}

// TestOperationAddLiquiditySuiFungibleToken is an operational test to add liquidity in the MUSE/SuiFungibleToken pool
func TestOperationAddLiquiditySuiFungibleToken(r *runner.E2ERunner, args []string) {
	liqMUSE, liqSuiToken := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySuiFungibleToken(liqMUSE, liqSuiToken)
}

// TestOperationAddLiquidityTON is an operational test to add liquidity in the MUSE/TON pool
func TestOperationAddLiquidityTON(r *runner.E2ERunner, args []string) {
	liqMUSE, liqTON := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityTON(liqMUSE, liqTON)
}
