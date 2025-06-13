package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v10 "github.com/RWAs-labs/muse/x/observer/migrations/v10"
	v11 "github.com/RWAs-labs/muse/x/observer/migrations/v11"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	observerKeeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		observerKeeper: keeper,
	}
}

// Migrate9to10 migrates the store from consensus version 9 to 10
func (m Migrator) Migrate9to10(ctx sdk.Context) error {
	return v10.MigrateStore(ctx, m.observerKeeper)
}

// Migrate10to11 migrates the store from consensus version 10 to 11
func (m Migrator) Migrate10to11(ctx sdk.Context) error {
	return v11.MigrateStore(ctx, m.observerKeeper)
}
