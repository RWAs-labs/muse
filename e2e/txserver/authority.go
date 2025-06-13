package txserver

import (
	"fmt"

	e2eutils "github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
)

// AddAuthorization adds a new authorization in the authority module for admin message
func (zts *MuseTxServer) AddAuthorization(msgURL string) error {
	// retrieve account
	accAdmin, err := zts.clientCtx.Keyring.Key(e2eutils.AdminPolicyName)
	if err != nil {
		return err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return err
	}

	// add new authorization
	_, err = zts.BroadcastTx(e2eutils.AdminPolicyName, authoritytypes.NewMsgAddAuthorization(
		addrAdmin.String(),
		msgURL,
		authoritytypes.PolicyType_groupAdmin,
	))
	if err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	return nil
}

// UpdateChainInfo sets the chain info in the authority module
func (zts MuseTxServer) UpdateChainInfo(chain chains.Chain) error {
	// retrieve account
	accAdmin, err := zts.clientCtx.Keyring.Key(e2eutils.AdminPolicyName)
	if err != nil {
		return err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return err
	}

	// set chain info
	_, err = zts.BroadcastTx(e2eutils.AdminPolicyName, authoritytypes.NewMsgUpdateChainInfo(
		addrAdmin.String(),
		chain,
	))
	if err != nil {
		return fmt.Errorf("failed to update chain info: %w", err)
	}

	return nil
}
