package keeper

import (
	"context"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdutils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateResource(goCtx context.Context, msg *types.MsgCreateResource) (*types.MsgCreateResourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate corresponding DIDDoc exists
	namespace := k.cheqdKeeper.GetDidNamespace(ctx)
	did := cheqdutils.JoinDID(cheqdtypes.DidMethod, namespace, msg.Payload.CollectionId)
	if !k.cheqdKeeper.HasDid(&ctx, did) {
		return nil, cheqdtypes.ErrDidDocNotFound.Wrapf(did)
	}

	// Validate Resource doesn't exist
	if k.HasResource(&ctx, msg.Payload.CollectionId, msg.Payload.Id) {
		return nil, types.ErrResourceExists.Wrap(msg.Payload.Id)
	}








	// getDid
	// get signatures for did modification
	//

	// Verify signatures
	signers := GetSignerDIDsForResourceCreation(resource)
	for _, signer := range signers {
		signature, found := cheqdtypes.FindSignInfoBySigner(msg.Signatures, signer)

		if !found {
			return nil, cheqdtypes.ErrSignatureNotFound.Wrapf("signer: %s", signer)
		}

		err := VerifySignature(&k.Keeper, &ctx, inMemoryResources, msg.Payload.GetSignBytes(), signature)
		if err != nil {
			return nil, err
		}
	}

	// Build Resource
	resource := msg.Payload.ToResource()
	// set created, checksum

	// Apply changes
	err := k.AppendResource(&ctx, &resource)
	if err != nil {
		return nil, types.ErrInternal.Wrapf(err.Error())
	}

	// Build and return response
	return &types.MsgCreateResourceResponse{
		Resource: &resource,
	}, nil
}

func GetSignerDIDsForResourceCreation(resource types.Resource) []string {
	//TODO: implement
	return []string{}
}
