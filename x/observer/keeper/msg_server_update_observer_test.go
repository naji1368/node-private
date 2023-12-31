package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_UpdateObserver(t *testing.T) {
	t.Run("successfully update tombstoned observer", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		assert.NoError(t, err)
		acc, found := k.GetNodeAccount(ctx, newOperatorAddress.String())
		assert.True(t, found)
		assert.Equal(t, newOperatorAddress.String(), acc.Operator)
	})

	t.Run("unable to update to a non validator address", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		assert.ErrorIs(t, err, types.ErrUpdateObserver)
	})

	t.Run("unable to update tombstoned validator with with non operator account", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)
		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            sample.AccAddress(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		assert.ErrorIs(t, err, types.ErrUpdateObserver)
	})
	t.Run("unable to update non-tombstoned observer with update reason tombstoned", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)
		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		assert.ErrorIs(t, err, types.ErrUpdateObserver)
	})
	t.Run("unable to update observer with no node account", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)
		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		assert.ErrorIs(t, err, types.ErrNodeAccountNotFound)
	})
	t.Run("unable to update observer when last observer count is missing", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)
		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            accAddressOfValidator.String(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_Tombstoned,
		})
		assert.ErrorIs(t, err, types.ErrLastObserverCountNotFound)
	})
	t.Run("update observer using admin policy", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)
		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)
		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            admin,
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
		})
		assert.NoError(t, err)
		acc, found := k.GetNodeAccount(ctx, newOperatorAddress.String())
		assert.True(t, found)
		assert.Equal(t, newOperatorAddress.String(), acc.Operator)
	})
	t.Run("fail to update observer using regular account and update type admin", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// #nosec G404 test purpose - weak randomness is not an issue here
		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validatorNew := sample.Validator(t, r)
		validatorNew.Status = stakingtypes.Bonded
		k.GetStakingKeeper().SetValidator(ctx, validatorNew)
		k.GetStakingKeeper().SetValidator(ctx, validator)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		chains := k.GetParams(ctx).GetSupportedChains()
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		assert.NoError(t, err)
		count := uint64(0)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
			count += 1
		}

		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator: accAddressOfValidator.String(),
		})

		newOperatorAddress, err := types.GetAccAddressFromOperatorAddress(validatorNew.OperatorAddress)
		assert.NoError(t, err)

		k.SetLastObserverCount(ctx, &types.LastObserverCount{
			Count: count,
		})

		_, err = srv.UpdateObserver(sdk.WrapSDKContext(ctx), &types.MsgUpdateObserver{
			Creator:            sample.AccAddress(),
			OldObserverAddress: accAddressOfValidator.String(),
			NewObserverAddress: newOperatorAddress.String(),
			UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
		})
		assert.ErrorIs(t, err, types.ErrUpdateObserver)
	})
}

func TestUpdateObserverList(t *testing.T) {
	t.Run("update observer list", func(t *testing.T) {
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		list := []string{sample.AccAddress(), sample.AccAddress(), sample.AccAddress(), oldObserverAddress}
		assert.Equal(t, oldObserverAddress, list[3])
		keeper.UpdateObserverList(list, oldObserverAddress, newObserverAddress)
		assert.Equal(t, 4, len(list))
		assert.Equal(t, newObserverAddress, list[3])
	})
}

func TestKeeper_UpdateObserverAddress(t *testing.T) {
	t.Run("update observer address", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		chains := k.GetParams(ctx).GetSupportedChains()
		observerList := CreateRandomObserverList(10, oldObserverAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  observerList,
			})
		}
		k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		observerMappers := k.GetAllObserverMappers(ctx)
		for _, om := range observerMappers {
			assert.Equal(t, len(observerList), len(om.ObserverList))
			assert.Equal(t, newObserverAddress, om.ObserverList[len(om.ObserverList)-1])
		}
	})
	t.Run("update observer address long observerList", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		chains := k.GetParams(ctx).GetSupportedChains()
		observerList := CreateRandomObserverList(1000, oldObserverAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  observerList,
			})
		}
		k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		observerMappers := k.GetAllObserverMappers(ctx)
		for _, om := range observerMappers {
			assert.Equal(t, len(observerList), len(om.ObserverList))
			assert.Equal(t, newObserverAddress, om.ObserverList[len(om.ObserverList)-1])
		}
	})
	t.Run("update observer address super long observerList", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		chains := k.GetParams(ctx).GetSupportedChains()
		observerList := CreateRandomObserverList(100000, oldObserverAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  observerList,
			})
		}
		k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		observerMappers := k.GetAllObserverMappers(ctx)
		for _, om := range observerMappers {
			assert.Equal(t, len(observerList), len(om.ObserverList))
			assert.Equal(t, newObserverAddress, om.ObserverList[len(om.ObserverList)-1])
		}
	})
	t.Run("update observer address short observerList", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		oldObserverAddress := sample.AccAddress()
		newObserverAddress := sample.AccAddress()
		chains := k.GetParams(ctx).GetSupportedChains()
		observerList := CreateRandomObserverList(1, oldObserverAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  observerList,
			})
		}
		k.UpdateObserverAddress(ctx, oldObserverAddress, newObserverAddress)
		observerMappers := k.GetAllObserverMappers(ctx)
		for _, om := range observerMappers {
			assert.Equal(t, len(observerList), len(om.ObserverList))
			assert.Equal(t, newObserverAddress, om.ObserverList[len(om.ObserverList)-1])
		}
	})
}

func CreateRandomObserverList(maxLen int, observerAddress string) []string {
	r := rand.New(rand.NewSource(9))
	list := make([]string, r.Intn(maxLen)+1)
	for i := range list {
		list[i] = sample.AccAddress()
	}
	list = append(list, observerAddress)
	return list
}
