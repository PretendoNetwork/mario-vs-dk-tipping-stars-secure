package nex_datastore

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func ChangeMeta(err error, client *nex.Client, callID uint32, dataStoreChangeMetaParam *nexproto.DataStoreChangeMetaParam) {
	// TODO - Check error and respond appropriately
	_ = database.UpdateMetaBinaryByDataStoreChangeMetaParam(dataStoreChangeMetaParam)

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodChangeMeta, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
