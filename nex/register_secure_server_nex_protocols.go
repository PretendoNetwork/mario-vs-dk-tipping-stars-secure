package nex

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	common_datastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
)

func registerSecureServerProtocols() {
	dataStoreProtocol := datastore.NewProtocol()
	commonDataStoreProtocol := common_datastore.NewCommonProtocol(dataStoreProtocol)
	globals.SecureEndpoint.RegisterServiceProtocol(dataStoreProtocol)

	commonDataStoreProtocol.SetManager(globals.DataStoreManager)
}
