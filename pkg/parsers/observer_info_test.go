package parsers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/RWAs-labs/muse/pkg/sdkconfig/default"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func TestParsefileToObserverMapper(t *testing.T) {
	file := "tmp.json"
	defer func(t *testing.T, fp string) {
		err := os.RemoveAll(fp)
		require.NoError(t, err)
	}(t, file)

	observerAddress := sample.AccAddress()
	commonGrantAddress := sample.AccAddress()
	validatorAddress := sample.AccAddress()

	createObserverList(file, observerAddress, commonGrantAddress, validatorAddress)
	obsListReadFromFile, err := ParsefileToObserverDetails(file)
	require.NoError(t, err)
	for _, obs := range obsListReadFromFile {
		require.Equal(
			t,
			obs.ObserverAddress,
			observerAddress,
		)
		require.Equal(
			t,
			obs.MuseClientGranteeAddress,
			commonGrantAddress,
		)
	}
}

func createObserverList(fp string, observerAddress, commonGrantAddress, validatorAddress string) {
	var listReader []ObserverInfoReader
	info := ObserverInfoReader{
		ObserverAddress:           observerAddress,
		MuseClientGranteeAddress:  commonGrantAddress,
		StakingGranteeAddress:     commonGrantAddress,
		StakingMaxTokens:          "100000000",
		StakingValidatorAllowList: []string{validatorAddress},
		SpendMaxTokens:            "100000000",
		GovGranteeAddress:         commonGrantAddress,
		MuseClientGranteePubKey:   "musepub1addwnpepqggtjvkmj6apcqr6ynyc5edxf2mpf5fxp2d3kwupemxtfwvg6gm7qv79fw0",
	}
	listReader = append(listReader, info)

	file, _ := json.MarshalIndent(listReader, "", " ")
	_ = ioutil.WriteFile(fp, file, 0600)
}
