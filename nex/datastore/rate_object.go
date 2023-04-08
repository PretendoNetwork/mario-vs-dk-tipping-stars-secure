package nex_datastore

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func RateObject(err error, client *nex.Client, callID uint32, target *datastore.DataStoreRatingTarget, param *datastore.DataStoreRateObjectParam, fetchRatings bool) {
	// TODO - Check error
	_ = database.UpdateRatingByDataIDAndSlot(uint32(target.DataID), target.Slot, param.RatingValue)

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	pRating := datastore.NewDataStoreRatingInfo()

	if fetchRatings {
		// TODO - Check error
		ratingInfo, _ := database.GetRatingByDataIDAndSlot(uint32(target.DataID), int(target.Slot))

		pRating = ratingInfo.Rating
	}

	rmcResponseStream.WriteStructure(pRating)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse.SetSuccess(datastore.MethodRateObject, rmcResponseBody)

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
