package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func register(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL) {
	localStation := stationUrls[0]
	publicStation := nex.NewStationURL("")

	// * Match the official server response format
	publicStation.SetScheme("prudp")
	publicStation.SetAddress(client.Address().IP.String())
	publicStation.SetPort(localStation.Port())
	publicStation.SetSID(localStation.SID())
	publicStation.SetType("3")
	publicStation.SetNatf("0")
	publicStation.SetNatm("0")
	publicStation.SetPmp("0")
	publicStation.SetUpnp("0")

	publicStationURL := publicStation.EncodeToString()

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt32LE(nex.Errors.Core.Unknown) // * Success
	rmcResponseStream.WriteUInt32LE(nexServer.ConnectionIDCounter().Increment())
	rmcResponseStream.WriteString(publicStationURL)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.SecureProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.SecureMethodRegister, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
