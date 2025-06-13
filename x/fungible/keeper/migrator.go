package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v3 "github.com/RWAs-labs/muse/x/fungible/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	fungibleKeeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		fungibleKeeper: keeper,
	}
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.fungibleKeeper)
}
