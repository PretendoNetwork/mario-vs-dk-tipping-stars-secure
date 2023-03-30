package nex_datastore

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/types"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func ChangeMeta(err error, client *nex.Client, callID uint32, dataStoreChangeMetaParam *nexproto.DataStoreChangeMetaParam) {
	// TODO - Check error and respond appropriately
	_ = database.UpdateMetaBinaryByDataStoreChangeMetaParam(dataStoreChangeMetaParam)

	if dataStoreChangeMetaParam.DataType == 123 {
		// TODO - Improve this! This just assumes the latest tip in the bucket is the only current tip! Could be VERY wrong!
		// * Tip Bucket update
		tipBucketStream := nex.NewStreamIn(dataStoreChangeMetaParam.MetaBinary, globals.NEXServer)

		tipBucket := types.NewTipBucket()
		tipBucket.ExtractFromStream(tipBucketStream)

		if len(tipBucket.Tips) >= 5 {
			//latestTip := tipBucket.Tips[len(tipBucket.Tips)-1]

			// TODO - Check error
			//_ = database.UpdateRatingByDataIDAndSlot(uint32(dataStoreChangeMetaParam.DataID), 0, int32(latestTip.Stars))
			//_ = database.UpdateRatingByDataIDAndSlot(uint32(dataStoreChangeMetaParam.DataID), 1, 1)
		}
	}

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
