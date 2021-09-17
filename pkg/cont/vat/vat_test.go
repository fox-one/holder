package vat

import (
	"context"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/mock"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/cont/sys"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/holder/pkg/number"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/holder/store/property"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestVault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pools := createPools(ctrl)
	vaults := createVaults(ctrl)
	properties := property.Memory()

	properties.Save(context.TODO(), sys.SystemProfitRateKey, "0.2")

	outputs := number.Values{}
	wallets := mock.NewMockWalletStore(ctrl)
	wallets.EXPECT().
		CreateTransfers(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, transfers []*core.Transfer) {
			for _, transfer := range transfers {
				outputs.Add(transfer.AssetID, transfer.Amount)
			}
		}).AnyTimes()

	assetz := mock.NewMockAssetService(ctrl)
	assetz.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&core.Asset{}, nil).AnyTimes()

	lockHandler := HandleLock(pools, vaults, assetz)
	releaseHandler := HandleRelease(pools, vaults, wallets, properties)

	assetID := newUUID()
	userID := newUUID()
	var version int64

	t.Run("lock 5 x 720", func(t *testing.T) {
		version += 1
		r := &cont.Request{
			Now:      time.Now(),
			Version:  version,
			TraceID:  newUUID(),
			Sender:   userID,
			FollowID: uuid.New(),
			AssetID:  assetID,
			Amount:   number.Decimal("5"),
			Action:   core.ActionVaultLock,
		}
		r.Body, _ = mtg.Encode(7200, 0)

		oldPool, _ := pools.Find(r.Context(), assetID)
		assert.Nil(t, lockHandler(r))
		newPool, _ := pools.Find(r.Context(), assetID)
		vault, _ := vaults.Find(r.Context(), r.TraceID)

		assert.True(t, vault.Liquidity.Equal(vault.Share))
		assert.True(t, r.Amount.Equal(vault.Amount))
		assert.True(t, newPool.Amount.Sub(oldPool.Amount).Equal(vault.Amount))
		assert.True(t, newPool.Share.Sub(oldPool.Share).Equal(vault.Share))
		assert.True(t, newPool.Liquidity.Sub(oldPool.Liquidity).Equal(vault.Liquidity))
	})

	t.Run("lock 10 x 720", func(t *testing.T) {
		version += 1
		r := &cont.Request{
			Now:      time.Now(),
			Version:  version,
			TraceID:  newUUID(),
			Sender:   userID,
			FollowID: uuid.New(),
			AssetID:  assetID,
			Amount:   number.Decimal("10"),
			Action:   core.ActionVaultLock,
		}
		r.Body, _ = mtg.Encode(7200, 0)

		oldPool, _ := pools.Find(r.Context(), assetID)
		assert.Nil(t, lockHandler(r))
		newPool, _ := pools.Find(r.Context(), assetID)
		vault, _ := vaults.Find(r.Context(), r.TraceID)

		assert.True(t, r.Amount.Equal(vault.Amount))
		assert.True(t, newPool.Amount.Sub(oldPool.Amount).Equal(vault.Amount))
		assert.True(t, newPool.Share.Sub(oldPool.Share).Equal(vault.Share))
		assert.True(t, newPool.Liquidity.Sub(oldPool.Liquidity).Equal(vault.Liquidity))
	})

	t.Run("lock & release", func(t *testing.T) {
		version += 1
		r := &cont.Request{
			Now:      time.Now(),
			Version:  version,
			TraceID:  newUUID(),
			Sender:   userID,
			FollowID: uuid.New(),
			AssetID:  assetID,
			Amount:   number.Decimal("8"),
			Action:   core.ActionVaultLock,
		}

		r.Body, _ = mtg.Encode(7200, 100)

		oldPool, _ := pools.Find(r.Context(), assetID)
		assert.Nil(t, lockHandler(r))
		newPool, _ := pools.Find(r.Context(), assetID)
		vault, _ := vaults.Find(r.Context(), r.TraceID)

		assert.Equal(t, int64(7200), vault.Duration)
		assert.Equal(t, int64(100), vault.MinDuration)
		assert.True(t, r.Amount.Equal(vault.Amount))
		assert.True(t, newPool.Amount.Sub(oldPool.Amount).Equal(vault.Amount))
		assert.True(t, newPool.Share.Sub(oldPool.Share).Equal(vault.Share))
		assert.True(t, newPool.Liquidity.Sub(oldPool.Liquidity).Equal(vault.Liquidity))

		t.Run("release before min duration", func(t *testing.T) {
			version += 1
			r := &cont.Request{
				Now:      time.Now(),
				Version:  version,
				TraceID:  newUUID(),
				Sender:   userID,
				FollowID: uuid.New(),
				AssetID:  assetID,
				Amount:   number.Decimal("8"),
				Action:   core.ActionVaultRelease,
			}

			r.Body, _ = mtg.Encode(types.UUID(vault.TraceID))

			err := releaseHandler(r)
			assert.NotNil(t, err)
			assert.Equal(t, "Vat/not-due", err.Error())
		})

		t.Run("release after min duration but before duration", func(t *testing.T) {
			version += 1
			r := &cont.Request{
				Now:      vault.CreatedAt.Add(time.Duration(vault.Duration/2) * time.Second),
				Version:  version,
				TraceID:  newUUID(),
				Sender:   userID,
				FollowID: uuid.New(),
				AssetID:  assetID,
				Amount:   number.Decimal("8"),
				Action:   core.ActionVaultRelease,
			}

			r.Body, _ = mtg.Encode(types.UUID(vault.TraceID))
			assert.Nil(t, releaseHandler(r))

			vault, _ = vaults.Find(r.Context(), vault.TraceID)
			assert.Equal(t, decimal.Zero.String(), vault.Reward.String())
			assert.Equal(t, vault.Penalty.String(), vault.Amount.Div(decimal.NewFromInt(2)).String())

			pool, _ := pools.Find(r.Context(), assetID)
			assert.Equal(t, vault.Penalty.String(), pool.Reward.Add(pool.Profit).String())
		})
	})

	t.Run("lock new and should have no reward", func(t *testing.T) {
		version += 1
		r := &cont.Request{
			Now:      time.Now(),
			Version:  version,
			TraceID:  newUUID(),
			Sender:   userID,
			FollowID: uuid.New(),
			AssetID:  assetID,
			Amount:   number.Decimal("10"),
			Action:   core.ActionVaultLock,
		}
		r.Body, _ = mtg.Encode(7200, 0)

		assert.Nil(t, lockHandler(r))
		pool, _ := pools.Find(r.Context(), assetID)
		vault, _ := vaults.Find(r.Context(), r.TraceID)

		assert.Equal(t, decimal.Zero.String(), GetReward(pool, vault).String())
	})

	t.Run("release remains", func(t *testing.T) {
		userVaults, _ := vaults.ListUser(context.TODO(), userID)
		assert.NotEmpty(t, userVaults)

		var totalAmount decimal.Decimal

		for _, vault := range userVaults {
			totalAmount = totalAmount.Add(vault.Amount)

			if vault.Status == core.VaultStatusReleased {
				continue
			}

			version += 1
			r := &cont.Request{
				Now:      vault.CreatedAt.Add(time.Hour * 24),
				Version:  version,
				TraceID:  newUUID(),
				Sender:   userID,
				FollowID: uuid.New(),
				AssetID:  assetID,
				Amount:   number.Decimal("8"),
				Action:   core.ActionVaultRelease,
			}

			r.Body, _ = mtg.Encode(types.UUID(vault.TraceID))
			assert.Nil(t, releaseHandler(r))
		}

		pool, _ := pools.Find(context.TODO(), assetID)
		assert.Equal(t, "0", pool.Amount.String())
		assert.Equal(t, "0", pool.Share.String())
		assert.Equal(t, "0", pool.Liquidity.String())
		assert.Equal(t, "0", pool.Reward.String())

		assert.Equal(t, outputs.Get(assetID).String(), totalAmount.Sub(pool.Profit).String())
	})
}

func newUUID() string {
	return uuid.New()
}

func createPools(ctrl *gomock.Controller) core.PoolStore {
	pools := map[string][]byte{}
	store := mock.NewMockPoolStore(ctrl)

	store.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, pool *core.Pool, version int64) {
			pool.Version = version
			b, _ := json.Marshal(pool)
			pools[pool.ID] = b
		}).AnyTimes()

	store.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id string) (*core.Pool, error) {
			pool := core.Pool{ID: id}
			if b, ok := pools[id]; ok {
				_ = json.Unmarshal(b, &pool)
			}

			return &pool, nil
		}).AnyTimes()

	return store
}

func createVaults(ctrl *gomock.Controller) core.VaultStore {
	vaults := map[string][]byte{}
	store := mock.NewMockVaultStore(ctrl)

	store.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, vault *core.Vault) {
			vault.ID = rand.Int63()
			b, _ := json.Marshal(vault)
			vaults[vault.TraceID] = b
		}).AnyTimes()

	store.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id string) (*core.Vault, error) {
			vault := core.Vault{TraceID: id}
			if b, ok := vaults[id]; ok {
				_ = json.Unmarshal(b, &vault)
			}

			return &vault, nil
		}).AnyTimes()

	store.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, vault *core.Vault, version int64) {
			vault.Version = version
			b, _ := json.Marshal(vault)
			vaults[vault.TraceID] = b
		}).AnyTimes()

	store.EXPECT().
		ListUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, userID string) ([]*core.Vault, error) {
			var results []*core.Vault
			for _, b := range vaults {
				var vault core.Vault
				_ = json.Unmarshal(b, &vault)

				if vault.UserID == userID {
					results = append(results, &vault)
				}
			}

			return results, nil
		}).AnyTimes()

	return store
}
