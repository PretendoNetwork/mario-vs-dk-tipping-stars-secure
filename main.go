package main

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/nex-go"
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

var nexServer *nex.Server

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(30701)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	nexServer.SetAccessKey("d8927c3f")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==MvDK:TS - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("==================")
	})

	nexServer.On("Packet", func(packet *nex.PacketV1) {
		fmt.Println(packet.Type())
	})

	secureConnectionProtocol := secureconnection.NewCommonSecureConnectionProtocol(nexServer)

	secureConnectionProtocol.AddConnection(addConnection)             // * Stubbed
	secureConnectionProtocol.UpdateConnection(updateConnection)       // * Stubbed
	secureConnectionProtocol.DoesConnectionExist(doesConnectionExist) // * Stubbed

	secureConnectionProtocol.Register(register) // * Override the common handler becuase needs a specific format maybe?

	nexServer.Listen(":60041")
}
