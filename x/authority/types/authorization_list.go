package types

import (
	"fmt"

	"cosmossdk.io/errors"
)

var (
	// OperationPolicyMessages keeps track of the message URLs that can, by default, only be executed by operational policy address
	OperationPolicyMessages = []string{
		"/musechain.musecore.crosschain.MsgRefundAbortedCCTX",
		"/musechain.musecore.crosschain.MsgAbortStuckCCTX",
		"/musechain.musecore.crosschain.MsgUpdateRateLimiterFlags",
		"/musechain.musecore.fungible.MsgDeploySystemContracts",
		"/musechain.musecore.fungible.MsgUpdateMRC20LiquidityCap",
		"/musechain.musecore.fungible.MsgUpdateMRC20WithdrawFee",
		"/musechain.musecore.fungible.MsgUnpauseMRC20",
		"/musechain.musecore.observer.MsgResetChainNonces",
		"/musechain.musecore.observer.MsgEnableCCTX",
		"/musechain.musecore.observer.MsgUpdateGasPriceIncreaseFlags",
		"/musechain.musecore.observer.MsgUpdateOperationalFlags",
		"/musechain.musecore.observer.MsgUpdateOperationalChainParams",
	}
	// AdminPolicyMessages keeps track of the message URLs that can, by default, only be executed by admin policy address
	AdminPolicyMessages = []string{
		"/musechain.musecore.crosschain.MsgUpdateERC20CustodyPauseStatus",
		"/musechain.musecore.crosschain.MsgMigrateERC20CustodyFunds",
		"/musechain.musecore.crosschain.MsgMigrateTssFunds",
		"/musechain.musecore.crosschain.MsgUpdateTssAddress",
		"/musechain.musecore.crosschain.MsgWhitelistERC20",
		"/musechain.musecore.fungible.MsgUpdateContractBytecode",
		"/musechain.musecore.fungible.MsgUpdateSystemContract",
		"/musechain.musecore.fungible.MsgUpdateGatewayContract",
		"/musechain.musecore.fungible.MsgRemoveForeignCoin",
		"/musechain.musecore.fungible.MsgDeployFungibleCoinMRC20",
		"/musechain.musecore.fungible.MsgUpdateMRC20Name",
		"/musechain.musecore.observer.MsgUpdateObserver",
		"/musechain.musecore.observer.MsgAddObserver",
		"/musechain.musecore.observer.MsgRemoveChainParams",
		"/musechain.musecore.authority.MsgAddAuthorization",
		"/musechain.musecore.authority.MsgRemoveAuthorization",
		"/musechain.musecore.authority.MsgUpdateChainInfo",
		"/musechain.musecore.authority.MsgRemoveChainInfo",
		"/musechain.musecore.lightclient.MsgEnableHeaderVerification",
		"/musechain.musecore.observer.MsgUpdateChainParams",
	}
	// EmergencyPolicyMessages keeps track of the message URLs that can, by default, only be executed by emergency policy address
	EmergencyPolicyMessages = []string{
		"/musechain.musecore.crosschain.MsgAddInboundTracker",
		"/musechain.musecore.crosschain.MsgAddOutboundTracker",
		"/musechain.musecore.crosschain.MsgRemoveOutboundTracker",
		"/musechain.musecore.crosschain.MsgRemoveInboundTracker",
		"/musechain.musecore.fungible.MsgPauseMRC20",
		"/musechain.musecore.observer.MsgUpdateKeygen",
		"/musechain.musecore.observer.MsgDisableCCTX",
		"/musechain.musecore.observer.MsgDisableFastConfirmation",
		"/musechain.musecore.lightclient.MsgDisableHeaderVerification",
	}
)

// DefaultAuthorizationsList list is the list of authorizations that presently exist in the system.
// This is the minimum set of authorizations that are required to be set when the authorization table is deployed
func DefaultAuthorizationsList() AuthorizationList {
	authorizations := make(
		[]Authorization,
		len(OperationPolicyMessages)+len(AdminPolicyMessages)+len(EmergencyPolicyMessages),
	)
	index := 0
	for _, msgURL := range OperationPolicyMessages {
		authorizations[index] = Authorization{
			MsgUrl:           msgURL,
			AuthorizedPolicy: PolicyType_groupOperational,
		}
		index++
	}
	for _, msgURL := range AdminPolicyMessages {
		authorizations[index] = Authorization{
			MsgUrl:           msgURL,
			AuthorizedPolicy: PolicyType_groupAdmin,
		}
		index++
	}
	for _, msgURL := range EmergencyPolicyMessages {
		authorizations[index] = Authorization{
			MsgUrl:           msgURL,
			AuthorizedPolicy: PolicyType_groupEmergency,
		}
		index++
	}

	return AuthorizationList{
		Authorizations: authorizations,
	}
}

// SetAuthorization adds the authorization to the list. If the authorization already exists, it updates the policy.
func (a *AuthorizationList) SetAuthorization(authorization Authorization) {
	for i, auth := range a.Authorizations {
		if auth.MsgUrl == authorization.MsgUrl {
			a.Authorizations[i].AuthorizedPolicy = authorization.AuthorizedPolicy
			return
		}
	}
	a.Authorizations = append(a.Authorizations, authorization)
}

// RemoveAuthorization removes the authorization from the list. It should be called by the admin policy account.
func (a *AuthorizationList) RemoveAuthorization(msgURL string) {
	for i, auth := range a.Authorizations {
		if auth.MsgUrl == msgURL {
			a.Authorizations = append(a.Authorizations[:i], a.Authorizations[i+1:]...)
			return
		}
	}
}

// GetAuthorizedPolicy returns the policy for the given message url. If the message url is not found, it returns an error.
func (a *AuthorizationList) GetAuthorizedPolicy(msgURL string) (PolicyType, error) {
	for _, auth := range a.Authorizations {
		if auth.MsgUrl == msgURL {
			return auth.AuthorizedPolicy, nil
		}
	}
	return PolicyType_groupEmpty, ErrAuthorizationNotFound
}

// Validate checks if the authorization list is valid. It returns an error if the message url is duplicated with different policies.
// It does not check if the list is empty or not, as an empty list is also considered valid.
func (a *AuthorizationList) Validate() error {
	checkMsgUrls := make(map[string]bool)
	for _, authorization := range a.Authorizations {
		if checkMsgUrls[authorization.MsgUrl] {
			return errors.Wrap(
				ErrInvalidAuthorizationList,
				fmt.Sprintf("duplicate message url: %s", authorization.MsgUrl),
			)
		}
		checkMsgUrls[authorization.MsgUrl] = true
	}
	return nil
}
