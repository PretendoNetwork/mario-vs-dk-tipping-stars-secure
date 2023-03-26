package nex_datastore

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func GetMetas(err error, client *nex.Client, callID uint32, dataIDs []uint64, param *nexproto.DataStoreGetMetaParam) {
	pMetaInfos := make([]*nexproto.DataStoreMetaInfo, 0, len(dataIDs))
	pResults := make([]*nex.Result, 0, len(dataIDs))

	// * We have to loop over every one and not batch-query the data
	// * because the console expects results for every level *in order*
	// * even if the level doesn't exist and this is only way to do that
	// * right now. Batch-querying the data would potentially lose
	// * entries and would require additional checks to make sure the
	// * data is correct, negating any speed benefits
	for i := 0; i < len(dataIDs); i++ {
		dataID := dataIDs[i]

		metaBinary, err := database.GetMetaBinaryByDataID(uint32(dataID))

		pMetaInfo := nexproto.NewDataStoreMetaInfo()
		var pResult *nex.Result

		if err == sql.ErrNoRows {
			pResult = nex.NewResultError(nex.Errors.DataStore.NotFound)
		} else { // TODO - Check for more errors
			pMetaInfo.DataID = uint64(metaBinary.DataID)
			pMetaInfo.OwnerID = metaBinary.OwnerPID
			pMetaInfo.Size = 0
			pMetaInfo.Name = metaBinary.Name
			pMetaInfo.DataType = metaBinary.DataType
			pMetaInfo.MetaBinary = metaBinary.Buffer
			pMetaInfo.Permission = nexproto.NewDataStorePermission()
			pMetaInfo.Permission.Permission = metaBinary.Permission
			pMetaInfo.Permission.RecipientIds = make([]uint32, 0)
			pMetaInfo.DelPermission = nexproto.NewDataStorePermission()
			pMetaInfo.DelPermission.Permission = metaBinary.DeletePermission
			pMetaInfo.DelPermission.RecipientIds = make([]uint32, 0)
			pMetaInfo.CreatedTime = nex.NewDateTime(nex.NewDateTime(0).Now()) // TODO - Change this!!
			pMetaInfo.UpdatedTime = nex.NewDateTime(nex.NewDateTime(0).Now()) // TODO - Change this!!
			pMetaInfo.Period = metaBinary.Period
			pMetaInfo.Status = 0
			pMetaInfo.ReferredCnt = 0
			pMetaInfo.ReferDataID = 0
			pMetaInfo.Flag = metaBinary.Flag
			pMetaInfo.ReferredTime = nex.NewDateTime(nex.NewDateTime(0).Now()) // TODO - Change this!!
			pMetaInfo.ExpireTime = nex.NewDateTime(nex.NewDateTime(0).Now())   // TODO - Change this!!
			pMetaInfo.Tags = metaBinary.Tags
			pMetaInfo.Ratings = make([]*nexproto.DataStoreRatingInfoWithSlot, 0) // TODO - Store ratings in DB

			pResult = nex.NewResultSuccess(nex.Errors.Core.Unknown)
		}

		pMetaInfos = append(pMetaInfos, pMetaInfo)
		pResults = append(pResults, pResult)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pMetaInfos)
	rmcResponseStream.WriteListResult(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMetas, rmcResponseBody)

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