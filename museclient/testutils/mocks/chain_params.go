package mocks

import (
	"testing"

	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/museconnector.non-eth.sol"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/testutils"
	"github.com/RWAs-labs/muse/pkg/constant"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

func MockChainParams(chainID int64, confirmation uint64) observertypes.ChainParams {
	connectorAddr := constant.EVMZeroAddress
	if a, ok := testutils.ConnectorAddresses[chainID]; ok {
		connectorAddr = a.Hex()
	}

	erc20CustodyAddr := constant.EVMZeroAddress
	if a, ok := testutils.CustodyAddresses[chainID]; ok {
		erc20CustodyAddr = a.Hex()
	}

	gwAddress := ""
	if gw, ok := testutils.GatewayAddresses[chainID]; ok {
		gwAddress = gw
	}

	return observertypes.ChainParams{
		ChainId:                     chainID,
		ConfirmationCount:           confirmation, // it is deprecated still needed to by pass chain params validation
		MuseTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    connectorAddr,
		Erc20CustodyContractAddress: erc20CustodyAddr,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             1,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		GatewayAddress:              gwAddress,
		IsSupported:                 true,
		ConfirmationParams: &observertypes.ConfirmationParams{
			SafeInboundCount:  confirmation,
			SafeOutboundCount: confirmation,
		},
	}
}

func MockConnectorNonEth(t *testing.T, chainID int64) *museconnector.MuseConnectorNonEth {
	connector, err := museconnector.NewMuseConnectorNonEth(testutils.ConnectorAddresses[chainID], &ethclient.Client{})
	require.NoError(t, err)
	return connector
}

func MockERC20Custody(t *testing.T, chainID int64) *erc20custody.ERC20Custody {
	custody, err := erc20custody.NewERC20Custody(testutils.CustodyAddresses[chainID], &ethclient.Client{})
	require.NoError(t, err)
	return custody
}
