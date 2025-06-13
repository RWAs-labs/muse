package mocks

import (
	"errors"

	"github.com/stretchr/testify/mock"

	keyinterfaces "github.com/RWAs-labs/muse/museclient/keys/interfaces"
	"github.com/RWAs-labs/muse/pkg/chains"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

var errSomethingIsWrong = errors.New("oopsie")

// Note that this is NOT codegen but a handwritten mock improvement.

func (_m *MusecoreClient) WithKeys(keys keyinterfaces.ObserverKeys) *MusecoreClient {
	_m.On("GetKeys").Maybe().Return(keys)

	return _m
}

func (_m *MusecoreClient) WithMuseChain() *MusecoreClient {
	_m.On("Chain").Maybe().Return(chains.MuseChainMainnet)

	return _m
}

func (_m *MusecoreClient) WithPostVoteOutbound(museTxHash string, ballotIndex string) *MusecoreClient {
	_m.On("PostVoteOutbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(museTxHash, ballotIndex, nil)

	return _m
}

func (_m *MusecoreClient) WithPostOutboundTracker(museTxHash string) *MusecoreClient {
	on := _m.On("PostOutboundTracker", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	if museTxHash != "" {
		on.Return(museTxHash, nil)
	} else {
		on.Return("", errSomethingIsWrong)
	}

	return _m
}

func (_m *MusecoreClient) WithPostVoteInbound(museTxHash string, ballotIndex string) *MusecoreClient {
	_m.On("PostVoteInbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(museTxHash, ballotIndex, nil)

	return _m
}

func (_m *MusecoreClient) WithRateLimiterFlags(flags *crosschaintypes.RateLimiterFlags) *MusecoreClient {
	on := _m.On("GetRateLimiterFlags", mock.Anything).Maybe()
	if flags != nil {
		on.Return(*flags, nil)
	} else {
		on.Return(crosschaintypes.RateLimiterFlags{}, errSomethingIsWrong)
	}

	return _m
}

func (_m *MusecoreClient) MockGetCctxByHash(err error) *MusecoreClient {
	_m.On("GetCctxByHash", mock.Anything, mock.Anything).Return(nil, err)
	return _m
}

func (_m *MusecoreClient) MockGetBallotByID(ballotIndex string, err error) *MusecoreClient {
	_m.On("GetBallotByID", mock.Anything, ballotIndex).Return(&observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: ballotIndex,
		Voters:           nil,
		ObservationType:  observertypes.ObservationType_InboundTx,
		BallotStatus:     observertypes.BallotStatus_BallotInProgress,
	}, err)
	return _m
}

func (_m *MusecoreClient) WithRateLimiterInput(in *crosschaintypes.QueryRateLimiterInputResponse) *MusecoreClient {
	on := _m.On("GetRateLimiterInput", mock.Anything, mock.Anything).Maybe()
	if in != nil {
		on.Return(in, nil)
	} else {
		on.Return(nil, errSomethingIsWrong)
	}

	return _m
}

func (_m *MusecoreClient) WithPendingCctx(chainID int64, cctxs []*crosschaintypes.CrossChainTx) *MusecoreClient {
	totalPending := uint64(len(cctxs))

	_m.On("ListPendingCCTX", mock.Anything, chainID).Maybe().Return(cctxs, totalPending, nil)

	return _m
}
