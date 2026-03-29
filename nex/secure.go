package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = true

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 7, 1))
	globals.SecureServer.AccessKey = "d8927c3f"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== MvDK:TS - Secure ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("==================")
	})

	globals.MinIOManager = common_globals.NewMinIOManager(globals.MinIOClient)
	globals.DataStoreManager = common_globals.NewDataStoreManager(globals.SecureEndpoint, database.Postgres)
	globals.DataStoreManager.GetUserFriendPIDs = globals.GetUserFriendPIDs
	globals.DataStoreManager.SetS3Config(globals.S3Bucket, globals.S3Key, globals.MinIOManager)

	// * Register the common handlers first so that they can be overridden if needed
	registerCommonSecureServerProtocols()
	registerSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MVDK_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
