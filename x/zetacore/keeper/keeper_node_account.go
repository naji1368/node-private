package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetNodeAccount set a specific nodeAccount in the store from its index
func (k Keeper) SetNodeAccount(ctx sdk.Context, nodeAccount types.NodeAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))
	b := k.cdc.MustMarshal(&nodeAccount)
	store.Set(types.KeyPrefix(nodeAccount.Index), b)
}

// GetNodeAccount returns a nodeAccount from its index
func (k Keeper) GetNodeAccount(ctx sdk.Context, index string) (val types.NodeAccount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveNodeAccount removes a nodeAccount from the store
func (k Keeper) RemoveNodeAccount(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllNodeAccount returns all nodeAccount
func (k Keeper) GetAllNodeAccount(ctx sdk.Context) (list []types.NodeAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NodeAccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) NodeAccountAll(c context.Context, req *types.QueryAllNodeAccountRequest) (*types.QueryAllNodeAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var nodeAccounts []*types.NodeAccount
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	nodeAccountStore := prefix.NewStore(store, types.KeyPrefix(types.NodeAccountKey))

	pageRes, err := query.Paginate(nodeAccountStore, req.Pagination, func(key []byte, value []byte) error {
		var nodeAccount types.NodeAccount
		if err := k.cdc.Unmarshal(value, &nodeAccount); err != nil {
			return err
		}

		nodeAccounts = append(nodeAccounts, &nodeAccount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNodeAccountResponse{NodeAccount: nodeAccounts, Pagination: pageRes}, nil
}

func (k Keeper) NodeAccount(c context.Context, req *types.QueryGetNodeAccountRequest) (*types.QueryGetNodeAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetNodeAccount(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNodeAccountResponse{NodeAccount: &val}, nil
}

// MESSAGES

func (k msgServer) SetNodeKeys(goCtx context.Context, msg *types.MsgSetNodeKeys) (*types.MsgSetNodeKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("msg creator %s not valid", msg.Creator))
	}
	_, found := k.GetNodeAccount(ctx, msg.Creator)
	if !found {
		na := types.NodeAccount{
			Creator:     msg.Creator,
			Index:       msg.Creator,
			NodeAddress: addr,
			PubkeySet:   msg.PubkeySet,
			NodeStatus:  types.NodeStatus_Unknown,
		}
		k.SetNodeAccount(ctx, na)
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("msg creator %s already has a node account", msg.Creator))
	}

	return &types.MsgSetNodeKeysResponse{}, nil
}