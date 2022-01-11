package main

import (
	"encoding/hex"
	"io/ioutil"

	"github.com/coinbase/rosetta-cli/configuration"
	"github.com/coinbase/rosetta-sdk-go/storage/modules"
	"github.com/coinbase/rosetta-sdk-go/types"

	"github.com/oasisprotocol/oasis-core-rosetta-gateway/services"
	"github.com/oasisprotocol/oasis-core-rosetta-gateway/tests/common"
)

func getRosettaConfig(ni *types.NetworkIdentifier) *configuration.Configuration {
	// Create a configuration file for the local testnet.
	config := configuration.DefaultConfiguration()

	config.Network = ni

	config.DataDirectory = "/tmp/rosetta-cli-oasistests"

	testEntityAddress, testEntityKeyPair := common.TestEntity()
	config.Construction = &configuration.ConstructionConfiguration{
		PrefundedAccounts: []*modules.PrefundedAccount{
			{
				PrivateKeyHex: hex.EncodeToString(testEntityKeyPair.PrivateKey),
				AccountIdentifier: &types.AccountIdentifier{
					Address:    testEntityAddress,
					SubAccount: nil,
					Metadata:   nil,
				},
				CurveType: types.Edwards25519,
				Currency:  services.OasisCurrency,
			},
		},
		ConstructorDSLFile: "oasis.ros",
		EndConditions: map[string]int{
			"transfer": 42,
		},
	}
	dataHistoricalBalanceEnabled := true
	config.Data.HistoricalBalanceEnabled = &dataHistoricalBalanceEnabled
	dataEndTip := true
	config.Data.EndConditions = &configuration.DataEndConditions{
		Tip: &dataEndTip,
	}

	return config
}

func main() {
	var err error

	_, ni := common.NewRosettaClient()

	config := getRosettaConfig(ni)
	if err = ioutil.WriteFile("rosetta-cli-config.json", []byte(common.DumpJSON(config)), 0o600); err != nil {
		panic(err)
	}
}
