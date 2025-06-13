package keeper

import (
	"encoding/json"
	"log"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/RWAs-labs/muse/x/observer/types"
)

func EmitEventBallotCreated(ctx sdk.Context, ballot types.Ballot, observationHash, observationChain string) {
	err := ctx.EventManager().EmitTypedEvent(&types.EventBallotCreated{
		BallotIdentifier: ballot.BallotIdentifier,
		BallotType:       ballot.ObservationType.String(),
		ObservationHash:  observationHash,
		ObservationChain: observationChain,
	})
	if err != nil {
		ctx.Logger().Error("failed to emit EventBallotCreated : %s", err.Error())
	}
}

// vendor this code from github.com/coinbase/rosetta-sdk-go/types
func prettyPrintStruct(val interface{}) string {
	prettyStruct, err := json.MarshalIndent(
		val,
		"",
		" ",
	)
	if err != nil {
		log.Fatal(err)
	}

	return string(prettyStruct)
}

func EmitEventKeyGenBlockUpdated(ctx sdk.Context, keygen *types.Keygen) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventKeygenBlockUpdated{
		MsgTypeUrl:    sdk.MsgTypeURL(&types.MsgUpdateKeygen{}),
		KeygenBlock:   strconv.Itoa(int(keygen.BlockNumber)),
		KeygenPubkeys: prettyPrintStruct(keygen.GranteePubkeys),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventKeygenBlockUpdated :", err)
	}
}

func EmitEventAddObserver(
	ctx sdk.Context,
	observerCount uint64,
	operatorAddress, museclientGranteeAddress, museclientGranteePubkey string,
) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventNewObserverAdded{
		MsgTypeUrl:               sdk.MsgTypeURL(&types.MsgAddObserver{}),
		ObserverAddress:          operatorAddress,
		MuseclientGranteeAddress: museclientGranteeAddress,
		MuseclientGranteePubkey:  museclientGranteePubkey,
		ObserverLastBlockCount:   observerCount,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EmitEventAddObserver :", err)
	}
}
