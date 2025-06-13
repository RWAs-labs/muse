package rpcimportable

import (
	"testing"

	"github.com/RWAs-labs/muse/pkg/rpc"
	"github.com/RWAs-labs/muse/pkg/sdkconfig"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRPCImportable(t *testing.T) {
	_ = rpc.Clients{}
}

func TestCosmosSdkConfigUntouched(t *testing.T) {
	museCfg := sdk.NewConfig()
	sdkconfig.Set(museCfg, true)
	if museCfg.GetBech32AccountAddrPrefix() != sdkconfig.AccountAddressPrefix {
		t.Logf("museCfg account prefix is not %s", sdkconfig.AccountAddressPrefix)
		t.FailNow()
	}

	// ensure that importing/using muse sdkconfig does not mutate the global config
	globalConfig := sdk.GetConfig()
	if globalConfig.GetBech32AccountAddrPrefix() != sdk.Bech32MainPrefix {
		t.Logf("globalConfig account prefix is not %s", sdk.Bech32MainPrefix)
		t.FailNow()
	}
}
