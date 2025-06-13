package local

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const (
	StakeAmount = "1000000000000000000000amuse"
)

// addNewObserver is a test function that adds a new observer to the network.
func addNewObserver(r *runner.E2ERunner) {
	fundHotkeyAccountForNonValidatorNode(r)
	stakeToBecomeValidator(r)
	addNodeAccount(r)
	addObserverAccount(r)
	addGrants(r)
}

// stakeToBecomeValidator is a helper function that stakes tokens to become a validator.
func stakeToBecomeValidator(r *runner.E2ERunner) {
	stakeTokens, err := sdk.ParseCoinNormalized(StakeAmount)
	require.NoError(r, err, "failed to parse coin")

	validatorsKeyring := r.MuseTxServer.GetValidatorsKeyring()
	pubkey, err := utils.FetchNodePubkey("musecore-new-validator")
	require.NoError(r, err)
	var pk cryptotypes.PubKey
	cdc := r.MuseTxServer.GetCodec()
	err = cdc.UnmarshalInterfaceJSON([]byte(pubkey), &pk)
	require.NoError(r, err, "failed to unmarshal pubkey")

	operator := getNewValidatorInfo(r)
	address, err := operator.GetAddress()
	require.NoError(r, err, "failed to get address")

	msg, err := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(address).String(),
		pk,
		stakeTokens,
		stakingtypes.NewDescription("musecore-new-validator", "", "", "", ""),
		stakingtypes.NewCommissionRates(
			sdkmath.LegacyMustNewDecFromStr("0.10"),
			sdkmath.LegacyMustNewDecFromStr("0.20"),
			sdkmath.LegacyMustNewDecFromStr("0.010"),
		),
		sdkmath.OneInt(),
	)
	require.NoError(r, err, "failed to create MsgCreateValidator")

	museTxServer := r.MuseTxServer
	validatorsTxServer := museTxServer.UpdateKeyring(validatorsKeyring)

	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// addNodeAccount adds the node account of a new validator to the network.
func addNodeAccount(r *runner.E2ERunner) {
	observerInfo, err := utils.FetchHotkeyAddress("museclient-new-validator")
	require.NoError(r, err)
	msg := observertypes.MsgAddObserver{
		Creator:                 r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		ObserverAddress:         observerInfo.ObserverAddress,
		MuseclientGranteePubkey: observerInfo.MuseClientGranteePubKey,
		AddNodeAccountOnly:      true,
	}
	_, err = r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, &msg)
	require.NoError(r, err)
}

// addObserverAccount adds a validator account to the observer set.
func addObserverAccount(r *runner.E2ERunner) {
	observerInfo, err := utils.FetchHotkeyAddress("museclient-new-validator")
	require.NoError(r, err)
	msg := observertypes.MsgAddObserver{
		Creator:                 r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		ObserverAddress:         observerInfo.ObserverAddress,
		MuseclientGranteePubkey: observerInfo.MuseClientGranteePubKey,
		AddNodeAccountOnly:      false,
	}
	_, err = r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, &msg)
	require.NoError(r, err)
}

// addGrants adds the necessary grants between operator and hotkey accounts.
func addGrants(r *runner.E2ERunner) {
	observerInfo, err := utils.FetchHotkeyAddress("museclient-new-validator")
	require.NoError(r, err)
	txTypes := crosschaintypes.GetAllAuthzMuseclientTxTypes()
	validatorsKeyring := r.MuseTxServer.GetValidatorsKeyring()

	museTxServer := r.MuseTxServer
	validatorsTxServer := museTxServer.UpdateKeyring(validatorsKeyring)

	operator := getNewValidatorInfo(r)
	for _, txType := range txTypes {
		msg, err := authz.NewMsgGrant(
			sdk.MustAccAddressFromBech32(observerInfo.ObserverAddress),
			sdk.MustAccAddressFromBech32(observerInfo.MuseClientGranteeAddress),
			authz.NewGenericAuthorization(txType),
			nil,
		)
		require.NoError(r, err)
		_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
		require.NoError(r, err, "failed to broadcast transaction")
	}

	msg, err := authz.NewMsgGrant(
		sdk.MustAccAddressFromBech32(observerInfo.ObserverAddress),
		sdk.MustAccAddressFromBech32(r.MuseTxServer.MustGetAccountAddressFromName(utils.UserEmissionsWithdrawName)),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&emissionstypes.MsgWithdrawEmission{})),
		nil,
	)
	require.NoError(r, err)
	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// fundHotkeyAccountForNonValidatorNode funds the hotkey address of a new validator.
func fundHotkeyAccountForNonValidatorNode(r *runner.E2ERunner) {
	amount, err := sdk.ParseCoinNormalized("100000000000000000000amuse")
	require.NoError(r, err, "failed to parse coin")

	observerInfo, err := utils.FetchHotkeyAddress("museclient-new-validator")
	require.NoError(r, err)

	validatorsKeyring := r.MuseTxServer.GetValidatorsKeyring()
	museTxServer := r.MuseTxServer
	validatorsTxServer := museTxServer.UpdateKeyring(validatorsKeyring)

	operator := getNewValidatorInfo(r)
	operatorAddress, err := operator.GetAddress()
	require.NoError(r, err, "failed to get address for operator")

	msg := banktypes.MsgSend{
		FromAddress: operatorAddress.String(),
		ToAddress:   observerInfo.MuseClientGranteeAddress,
		Amount:      sdk.NewCoins(amount),
	}

	_, err = validatorsTxServer.BroadcastTx(operator.Name, &msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// getNewValidatorInfo retrieves the keyring record for the new validator from the keyring.
func getNewValidatorInfo(r *runner.E2ERunner) *keyring.Record {
	record, err := r.MuseTxServer.GetValidatorsKeyring().Key("operator-new-validator")
	require.NoError(r, err, "failed to get operator-new-validator key")
	return record
}
