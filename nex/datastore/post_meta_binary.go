package nex_datastore

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func PostMetaBinary(err error, client *nex.Client, callID uint32, dataStorePreparePostParam *datastore.DataStorePreparePostParam) {
	metaBinary := database.GetMetaBinaryByTypeAndOwnerPIDAndSlotID(dataStorePreparePostParam.DataType, client.PID(), uint8(dataStorePreparePostParam.PersistenceInitParam.PersistenceSlotId))

	if metaBinary.DataID != 0 {
		// * Meta binary already exists
		if dataStorePreparePostParam.PersistenceInitParam.DeleteLastObject {
			// * Delete existing object before uploading new one
			// TODO - Check error
			_ = database.DeleteMetaBinaryByDataID(metaBinary.DataID)
			// TODO - Delete old ratings?
		}
	}

	// TODO - See if this is actually always the case?
	// * Always upload a new object (?)
	dataID := database.InsertMetaBinaryByDataStorePreparePostParamWithOwnerPID(dataStorePreparePostParam, client.PID())

	for i := 0; i < len(dataStorePreparePostParam.RatingInitParams); i++ {
		ratingInitParam := dataStorePreparePostParam.RatingInitParams[i]

		// TODO - Check error
		_ = database.InsertRatingByDataIDAndDataStoreRatingInitParamWithSlot(dataID, ratingInitParam)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteUInt64LE(uint64(dataID))

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPostMetaBinary, rmcResponseBody)

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
