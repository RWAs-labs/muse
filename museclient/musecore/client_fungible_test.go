package musecore

import (
	"context"
	"testing"

	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/testutil/sample"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

func Test_GetForeignCoinsFromAsset(t *testing.T) {
	erc20Asset := sample.EthAddress()

	tests := []struct {
		name         string
		chainID      int64
		assetAddress ethcommon.Address
		errMsg       string
	}{
		{
			name:         "get ERC20 foreign coins from asset",
			chainID:      1,
			assetAddress: erc20Asset,
		},
		{
			name:         "get Gas foreign coins from zero-address asset",
			chainID:      1,
			assetAddress: ethcommon.Address{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// construct foreign coin
			assetString := tt.assetAddress.Hex()
			if crypto.IsEmptyAddress(tt.assetAddress) {
				assetString = ""
			}
			fCoins := sample.ForeignCoins(t, "0x123")
			fCoins.Asset = assetString
			fCoins.ForeignChainId = tt.chainID

			// mock RPC server
			method := "/musechain.musecore.fungible.Query/ForeignCoinsFromAsset"
			mockRequest := fungibletypes.QueryGetForeignCoinsFromAssetRequest{
				ChainId: fCoins.ForeignChainId,
				Asset:   fCoins.Asset,
			}
			mockResponse := &fungibletypes.QueryGetForeignCoinsFromAssetResponse{
				ForeignCoins: fCoins,
			}
			setupMockServer(t, fungibletypes.RegisterQueryServer, method, mockRequest, mockResponse)
			client := setupMusecoreClient(
				t,
				withDefaultObserverKeys(),
				withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
			)

			// ACT
			ctx := context.Background()
			resp, err := client.GetForeignCoinsFromAsset(ctx, tt.chainID, tt.assetAddress)

			// ASSERT
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Equal(t, fungibletypes.ForeignCoins{}, resp)
				return
			}
			require.NoError(t, err)
			require.Equal(t, fCoins, resp)
		})
	}
}
