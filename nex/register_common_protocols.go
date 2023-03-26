package nex

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	nex_secure_connection_common "github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/nex/secure_connection/common"
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

func registerCommonProtocols() {
	secureConnectionProtocol := secureconnection.NewCommonSecureConnectionProtocol(globals.NEXServer)

	secureConnectionProtocol.AddConnection(nex_secure_connection_common.AddConnection)             // * Stubbed
	secureConnectionProtocol.UpdateConnection(nex_secure_connection_common.UpdateConnection)       // * Stubbed
	secureConnectionProtocol.DoesConnectionExist(nex_secure_connection_common.DoesConnectionExist) // * Stubbed
}
