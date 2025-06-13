package signer

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/sui/client"
	musesui "github.com/RWAs-labs/muse/pkg/contracts/sui"
)

// withdrawAndCallObjRefs contains all the object references needed for withdraw and call
type withdrawAndCallObjRefs struct {
	gateway     sui.ObjectRef
	withdrawCap sui.ObjectRef
	onCall      []sui.ObjectRef
	suiCoins    []*sui.ObjectRef
}

// withdrawAndCallPTBArgs contains all the arguments needed for withdraw and call
type withdrawAndCallPTBArgs struct {
	withdrawAndCallObjRefs
	coinType  string
	amount    uint64
	nonce     uint64
	gasBudget uint64
	receiver  string
	payload   musesui.CallPayload
}

// withdrawAndCallPTB builds unsigned withdraw and call PTB Sui transaction
// it chains the following calls:
// 1. withdraw_impl on gateway
// 2. gas budget coin transfer to TSS
// 3. on_call on target contract
// The function returns a TxnMetaData object with tx bytes, the other fields are ignored
func (s *Signer) withdrawAndCallPTB(args withdrawAndCallPTBArgs) (tx models.TxnMetaData, err error) {
	var (
		tssAddress       = s.TSS().PubKey().AddressSui()
		gatewayPackageID = s.gateway.PackageID()
		gatewayModule    = s.gateway.Module()
		ptb              = suiptb.NewTransactionDataTransactionBuilder()
	)

	// Parse signer address
	signerAddr, err := sui.AddressFromHex(tssAddress)
	if err != nil {
		return tx, errors.Wrapf(err, "invalid signer address %s", tssAddress)
	}

	// Add withdraw_impl command and get its command index
	if err := ptbAddCmdWithdrawImpl(
		ptb,
		gatewayPackageID,
		gatewayModule,
		args.gateway,
		args.withdrawCap,
		args.coinType,
		args.amount,
		args.nonce,
		args.gasBudget,
	); err != nil {
		return tx, err
	}

	// Create arguments to access the two results from the withdraw_impl call
	cmdIndex := uint16(0)
	argWithdrawnCoins := suiptb.Argument{
		NestedResult: &suiptb.NestedResult{
			Cmd:    cmdIndex,
			Result: 0, // First result (main coins)
		},
	}

	argBudgetCoins := suiptb.Argument{
		NestedResult: &suiptb.NestedResult{
			Cmd:    cmdIndex,
			Result: 1, // Second result (budget coins)
		},
	}

	// Add gas budget transfer command
	err = ptbAddCmdGasBudgetTransfer(ptb, argBudgetCoins, *signerAddr)
	if err != nil {
		return tx, err
	}

	// Add on_call command
	err = ptbAddCmdOnCall(
		ptb,
		args.receiver,
		args.coinType,
		argWithdrawnCoins,
		args.onCall,
		args.payload,
	)
	if err != nil {
		return tx, err
	}

	// Finish building the PTB
	pt := ptb.Finish()

	// Wrap the PTB into a transaction data
	txData := suiptb.NewTransactionData(
		signerAddr,
		pt,
		args.suiCoins,
		args.gasBudget,
		suiclient.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		return tx, errors.Wrapf(err, "failed to marshal transaction data: %v", txData)
	}

	// Encode the transaction bytes to base64
	return models.TxnMetaData{
		TxBytes: base64.StdEncoding.EncodeToString(txBytes),
	}, nil
}

// getWithdrawAndCallObjectRefs returns the SUI object references for withdraw and call
//   - Initial shared version will be used for shared objects
//   - Current version will be used for non-shared objects, e.g. withdraw cap
func (s *Signer) getWithdrawAndCallObjectRefs(
	ctx context.Context,
	withdrawCapID string,
	onCallObjectIDs []string,
	gasBudget uint64,
) (withdrawAndCallObjRefs, error) {
	objectIDs := append([]string{s.gateway.ObjectID(), withdrawCapID}, onCallObjectIDs...)

	// query objects in batch
	suiObjects, err := s.client.SuiMultiGetObjects(ctx, models.SuiMultiGetObjectsRequest{
		ObjectIds: objectIDs,
		Options: models.SuiObjectDataOptions{
			// show owner info in order to retrieve object initial shared version
			ShowOwner: true,
		},
	})
	if err != nil {
		return withdrawAndCallObjRefs{}, errors.Wrapf(err, "failed to get objects for %v", objectIDs)
	}

	// should never mismatch, just a sanity check
	if len(suiObjects) != len(objectIDs) {
		return withdrawAndCallObjRefs{}, fmt.Errorf("expected %d objects, but got %d", len(objectIDs), len(suiObjects))
	}

	// ensure no owned objects are used for on_call
	if err := client.CheckContainOwnedObject(suiObjects[2:]); err != nil {
		return withdrawAndCallObjRefs{}, errors.Wrapf(err, "objects used for on_call must be shared")
	}

	// convert object data to object references
	objectRefs := make([]sui.ObjectRef, len(objectIDs))

	for i, object := range suiObjects {
		objectID, err := sui.ObjectIdFromHex(object.Data.ObjectId)
		if err != nil {
			return withdrawAndCallObjRefs{}, errors.Wrapf(err, "failed to parse object ID %s", object.Data.ObjectId)
		}

		objectVersion, err := strconv.ParseUint(object.Data.Version, 10, 64)
		if err != nil {
			return withdrawAndCallObjRefs{}, errors.Wrapf(err, "failed to parse object version %s", object.Data.Version)
		}

		// must use initial version for shared object, not the current version
		// withdraw cap is not a shared object, so we must use current version
		if object.Data.ObjectId != withdrawCapID {
			objectVersion, err = musesui.ExtractInitialSharedVersion(*object.Data)
			if err != nil {
				return withdrawAndCallObjRefs{}, errors.Wrapf(
					err,
					"failed to extract initial shared version for object %s",
					object.Data.ObjectId,
				)
			}
		}

		objectDigest, err := sui.NewBase58(object.Data.Digest)
		if err != nil {
			return withdrawAndCallObjRefs{}, errors.Wrapf(err, "failed to parse object digest %s", object.Data.Digest)
		}

		objectRefs[i] = sui.ObjectRef{
			ObjectId: objectID,
			Version:  objectVersion,
			Digest:   objectDigest,
		}
	}

	// get latest TSS SUI coin object ref for gas payment
	suiCoinObjRefs, err := s.client.GetSuiCoinObjectRefs(ctx, s.TSS().PubKey().AddressSui(), gasBudget)
	if err != nil {
		return withdrawAndCallObjRefs{}, errors.Wrap(err, "unable to get TSS SUI coin objects")
	}

	return withdrawAndCallObjRefs{
		gateway:     objectRefs[0],
		withdrawCap: objectRefs[1],
		onCall:      objectRefs[2:],
		suiCoins:    suiCoinObjRefs,
	}, nil
}

// ptbAddCmdWithdrawImpl adds the withdraw_impl command to the PTB
func ptbAddCmdWithdrawImpl(
	ptb *suiptb.ProgrammableTransactionBuilder,
	gatewayPackageIDStr string,
	gatewayModule string,
	gatewayObjRef sui.ObjectRef,
	withdrawCapObjRef sui.ObjectRef,
	coinType string,
	amount uint64,
	nonce uint64,
	gasBudget uint64,
) error {
	// Parse gateway package ID
	gatewayPackageID, err := sui.PackageIdFromHex(gatewayPackageIDStr)
	if err != nil {
		return errors.Wrapf(err, "invalid gateway package ID %s", gatewayPackageIDStr)
	}

	// Parse coin type
	tagCoinType, err := musesui.TypeTagFromString(coinType)
	if err != nil {
		return errors.Wrapf(err, "invalid coin type %s", coinType)
	}

	// Create gateway object argument
	argGatewayObject, err := ptb.Obj(suiptb.ObjectArg{
		SharedObject: &suiptb.SharedObjectArg{
			Id:                   gatewayObjRef.ObjectId,
			InitialSharedVersion: gatewayObjRef.Version,
			Mutable:              true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "unable to create gateway object argument")
	}

	// Create amount argument
	argAmount, err := ptb.Pure(amount)
	if err != nil {
		return errors.Wrapf(err, "unable to create amount argument")
	}

	// Create nonce argument
	argNonce, err := ptb.Pure(nonce)
	if err != nil {
		return errors.Wrapf(err, "unable to create nonce argument")
	}

	// Create gas budget argument
	argGasBudget, err := ptb.Pure(gasBudget)
	if err != nil {
		return errors.Wrapf(err, "unable to create gas budget argument")
	}

	// Create withdraw cap argument
	argWithdrawCap, err := ptb.Obj(suiptb.ObjectArg{ImmOrOwnedObject: &withdrawCapObjRef})
	if err != nil {
		return errors.Wrapf(err, "unable to create withdraw cap object argument")
	}

	// add Move call for withdraw_impl
	// #nosec G115 always in range
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:  gatewayPackageID,
			Module:   gatewayModule,
			Function: musesui.FuncWithdrawImpl,
			TypeArguments: []sui.TypeTag{
				{Struct: &tagCoinType},
			},
			Arguments: []suiptb.Argument{
				argGatewayObject,
				argAmount,
				argNonce,
				argGasBudget,
				argWithdrawCap,
			},
		},
	})

	return nil
}

// ptbAddCmdGasBudgetTransfer adds the gas budget transfer command to the PTB
func ptbAddCmdGasBudgetTransfer(
	ptb *suiptb.ProgrammableTransactionBuilder,
	argBudgetCoins suiptb.Argument,
	signerAddr sui.Address,
) error {
	// create TSS address argument
	argTSSAddr, err := ptb.Pure(signerAddr)
	if err != nil {
		return errors.Wrapf(err, "unable to create tss address argument")
	}

	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{argBudgetCoins},
			Address: argTSSAddr,
		},
	})

	return nil
}

// ptbAddCmdOnCall adds the on_call command to the PTB
func ptbAddCmdOnCall(
	ptb *suiptb.ProgrammableTransactionBuilder,
	receiver string,
	coinTypeStr string,
	argWithdrawnCoins suiptb.Argument,
	onCallObjectRefs []sui.ObjectRef,
	cp musesui.CallPayload,
) error {
	// Parse target package ID
	targetPackageID, err := sui.PackageIdFromHex(receiver)
	if err != nil {
		return errors.Wrapf(err, "invalid target package ID %s", receiver)
	}

	// Parse coin type
	coinType, err := musesui.TypeTagFromString(coinTypeStr)
	if err != nil {
		return errors.Wrapf(err, "invalid coin type %s", coinTypeStr)
	}

	// Build the type arguments for on_call in order: [withdrawn coin type, ... payload type arguments]
	onCallTypeArgs := make([]sui.TypeTag, 0, len(cp.TypeArgs)+1)
	onCallTypeArgs = append(onCallTypeArgs, sui.TypeTag{Struct: &coinType})
	for _, typeArg := range cp.TypeArgs {
		typeStruct, err := musesui.TypeTagFromString(typeArg)
		if err != nil {
			return errors.Wrapf(err, "invalid type argument %s", typeArg)
		}
		onCallTypeArgs = append(onCallTypeArgs, sui.TypeTag{Struct: &typeStruct})
	}

	// Build the args for on_call: [withdrawns coins + payload objects + message]
	onCallArgs := make([]suiptb.Argument, 0, len(cp.ObjectIDs)+1)
	onCallArgs = append(onCallArgs, argWithdrawnCoins)

	// Add the payload objects, objects are all shared
	for _, onCallObjectRef := range onCallObjectRefs {
		objectArg, err := ptb.Obj(suiptb.ObjectArg{
			SharedObject: &suiptb.SharedObjectArg{
				Id:                   onCallObjectRef.ObjectId,
				InitialSharedVersion: onCallObjectRef.Version,
				Mutable:              true,
			},
		})
		if err != nil {
			return errors.Wrapf(err, "unable to create object argument: %v", onCallObjectRef)
		}
		onCallArgs = append(onCallArgs, objectArg)
	}

	// Add any additional message arguments
	messageArg, err := ptb.Pure(cp.Message)
	if err != nil {
		return errors.Wrapf(err, "unable to create message argument: %x", cp.Message)
	}
	onCallArgs = append(onCallArgs, messageArg)

	// Call the target contract on_call
	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       targetPackageID,
			Module:        musesui.ModuleConnected,
			Function:      musesui.FuncOnCall,
			TypeArguments: onCallTypeArgs,
			Arguments:     onCallArgs,
		},
	})

	return nil
}
