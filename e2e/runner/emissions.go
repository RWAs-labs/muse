package runner

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/e2e/txserver"
	e2eutils "github.com/RWAs-labs/muse/e2e/utils"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// FundEmissionsPool funds the emissions pool on MuseChain with the same value as used originally on mainnet (20M MUSE)
func (r *E2ERunner) FundEmissionsPool() error {
	r.Logger.Print("⚙️ funding the emissions pool on MuseChain with 20M MUSE (%s)", txserver.EmissionsPoolAddress)

	return r.MuseTxServer.FundEmissionsPool(e2eutils.OperationalPolicyName, EmissionsPoolFunding)
}

// WithdrawEmissions withdraws emissions from the emission pool on MuseChain for all observers
// This functions uses the UserEmissionsWithdrawName to create the withdraw tx.
// UserEmissionsWithdraw can sign the authz transactions because the necessary permissions are granted in the genesis file
func (r *E2ERunner) WithdrawEmissions() error {
	observerSet, err := r.ObserverClient.ObserverSet(r.Ctx, &observertypes.QueryObserverSet{})
	if err != nil {
		return err
	}

	for _, observer := range observerSet.Observers {
		r.Logger.Print("🏃 Withdrawing emissions for observer %s", observer)
		var (
			baseDenom            = config.BaseDenom
			queryObserverBalance = &banktypes.QueryBalanceRequest{
				Address: observer,
				Denom:   baseDenom,
			}
		)

		balanceBefore, err := r.BankClient.Balance(r.Ctx, queryObserverBalance)
		if err != nil {
			return errors.Wrapf(err, "failed to get balance for observer before withdrawing emissions %s", observer)
		}

		availableCoin, err := r.FetchWithdrawableEmissions(observer)
		if err != nil {
			return err
		}

		if availableCoin.Amount.IsZero() {
			r.Logger.Print("no emissions to withdraw for observer %s", observer)
			continue
		}

		if err := r.MuseTxServer.WithdrawAllEmissions(availableCoin.Amount, e2eutils.UserEmissionsWithdrawName, observer); err != nil {
			r.Logger.Error("failed to withdraw emissions for observer %s: %s", observer, err)
			r.Logger.Error("Withdraw amount: %s", availableCoin.Amount)
			availableCoinAfter, fetchErr := r.FetchWithdrawableEmissions(observer)
			if fetchErr != nil {
				r.Logger.Error("failed to fetch available emissions for observer %s: %s", observer, fetchErr)
				return err
			}
			r.Logger.Error("Available emissions after failed withdrawal: %s", availableCoinAfter.Amount)
			return err
		}

		balanceAfter, err := r.BankClient.Balance(r.Ctx, queryObserverBalance)
		if err != nil {
			return errors.Wrapf(err, "failed to get balance for observer after withdrawing emissions %s", observer)
		}

		changeInBalance := balanceAfter.Balance.Sub(*balanceBefore.Balance).Amount
		if !changeInBalance.Equal(availableCoin.Amount) {
			return fmt.Errorf(
				"invalid balance change for observer %s, expected %s, got %s",
				observer,
				availableCoin.Amount,
				changeInBalance,
			)
		}
	}

	return nil
}

func (r *E2ERunner) FetchWithdrawableEmissions(observer string) (sdk.Coin, error) {
	availableAmount, err := r.EmissionsClient.ShowAvailableEmissions(
		r.Ctx,
		&emissionstypes.QueryShowAvailableEmissionsRequest{
			Address: observer,
		},
	)
	if err != nil {
		return sdk.Coin{}, errors.Wrapf(err, "failed to get available emissions for observer %s", observer)
	}

	availableCoin, err := sdk.ParseCoinNormalized(availableAmount.Amount)
	if err != nil {
		return sdk.Coin{}, errors.Wrap(err, "failed to parse coin amount")
	}
	return availableCoin, nil
}
