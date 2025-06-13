package app

import (
	"context"
	"os"
	"path"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/stretchr/testify/require"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	lightclienttypes "github.com/RWAs-labs/muse/x/lightclient/types"
)

func TestUpgradeTracker(t *testing.T) {
	r := require.New(t)

	tmpdir := t.TempDir()

	allUpgrades := upgradeTracker{
		upgrades: []upgradeTrackerItem{
			{
				index: 1000,
				storeUpgrade: &storetypes.StoreUpgrades{
					Added: []string{authoritytypes.ModuleName},
				},
			},
			{
				index: 2000,
				storeUpgrade: &storetypes.StoreUpgrades{
					Added: []string{lightclienttypes.ModuleName},
				},
				upgradeHandler: func(ctx context.Context, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				},
			},
			{
				index: 3000,
				upgradeHandler: func(ctx context.Context, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				},
			},
		},
		stateFileDir: tmpdir,
	}

	upgradeHandlers, storeUpgrades := allUpgrades.mergeAllUpgrades()
	r.Len(storeUpgrades.Added, 2)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 0)
	r.Len(upgradeHandlers, 2)

	// should return all migrations on first call
	upgradeHandlers, storeUpgrades, err := allUpgrades.getIncrementalUpgrades()
	r.NoError(err)
	r.Len(storeUpgrades.Added, 2)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 0)
	r.Len(upgradeHandlers, 2)

	// should return no upgrades on second call
	upgradeHandlers, storeUpgrades, err = allUpgrades.getIncrementalUpgrades()
	r.NoError(err)
	r.Len(storeUpgrades.Added, 0)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 0)
	r.Len(upgradeHandlers, 0)

	// now add a upgrade and ensure that it gets run without running
	// the other upgrades
	allUpgrades.upgrades = append(allUpgrades.upgrades, upgradeTrackerItem{
		index: 4000,
		storeUpgrade: &storetypes.StoreUpgrades{
			Deleted: []string{"example"},
		},
	})

	upgradeHandlers, storeUpgrades, err = allUpgrades.getIncrementalUpgrades()
	r.NoError(err)
	r.Len(storeUpgrades.Added, 0)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 1)
	r.Len(upgradeHandlers, 0)
}

func TestUpgradeTrackerBadState(t *testing.T) {
	r := require.New(t)

	tmpdir := t.TempDir()

	stateFilePath := path.Join(tmpdir, incrementalUpgradeTrackerStateFile)

	err := os.WriteFile(stateFilePath, []byte("badstate"), 0o600)
	r.NoError(err)

	allUpgrades := upgradeTracker{
		upgrades:     []upgradeTrackerItem{},
		stateFileDir: tmpdir,
	}
	_, _, err = allUpgrades.getIncrementalUpgrades()
	r.Error(err)
}
