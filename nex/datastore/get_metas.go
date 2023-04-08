package nex_datastore

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func GetMetas(err error, client *nex.Client, callID uint32, dataIDs []uint64, param *datastore.DataStoreGetMetaParam) {
	pMetaInfos := make([]*datastore.DataStoreMetaInfo, 0, len(dataIDs))
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

		pMetaInfo := datastore.NewDataStoreMetaInfo()
		var pResult *nex.Result

		if err == sql.ErrNoRows {
			pMetaInfo.Permission = datastore.NewDataStorePermission()
			pMetaInfo.Permission.Permission = 0
			pMetaInfo.Permission.RecipientIds = make([]uint32, 0)
			pMetaInfo.DelPermission = datastore.NewDataStorePermission()
			pMetaInfo.DelPermission.Permission = 0
			pMetaInfo.DelPermission.RecipientIds = make([]uint32, 0)
			pMetaInfo.CreatedTime = nex.NewDateTime(0)
			pMetaInfo.UpdatedTime = nex.NewDateTime(0)
			pMetaInfo.ReferredTime = nex.NewDateTime(0)
			pMetaInfo.ExpireTime = nex.NewDateTime(0)
			pMetaInfo.Ratings = make([]*datastore.DataStoreRatingInfoWithSlot, 0)

			pResult = nex.NewResultError(nex.Errors.DataStore.NotFound)
		} else { // TODO - Check for more errors
			pMetaInfo.DataID = uint64(metaBinary.DataID)
			pMetaInfo.OwnerID = metaBinary.OwnerPID
			pMetaInfo.Size = 0
			pMetaInfo.Name = metaBinary.Name
			pMetaInfo.DataType = metaBinary.DataType
			pMetaInfo.MetaBinary = metaBinary.Buffer
			pMetaInfo.Permission = datastore.NewDataStorePermission()
			pMetaInfo.Permission.Permission = metaBinary.Permission
			pMetaInfo.Permission.RecipientIds = make([]uint32, 0)
			pMetaInfo.DelPermission = datastore.NewDataStorePermission()
			pMetaInfo.DelPermission.Permission = metaBinary.DeletePermission
			pMetaInfo.DelPermission.RecipientIds = make([]uint32, 0)
			pMetaInfo.CreatedTime = metaBinary.CreationTime
			pMetaInfo.UpdatedTime = metaBinary.UpdatedTime
			pMetaInfo.Period = metaBinary.Period
			pMetaInfo.Status = 0
			pMetaInfo.ReferredCnt = 0
			pMetaInfo.ReferDataID = 0
			pMetaInfo.Flag = metaBinary.Flag
			pMetaInfo.ReferredTime = metaBinary.ReferredTime
			pMetaInfo.ExpireTime = metaBinary.ExpireTime
			pMetaInfo.Tags = metaBinary.Tags
			pMetaInfo.Ratings = make([]*datastore.DataStoreRatingInfoWithSlot, 0) // TODO - Store ratings in DB

			pResult = nex.NewResultSuccess(nex.Errors.Core.Unknown)
		}

		pMetaInfos = append(pMetaInfos, pMetaInfo)
		pResults = append(pResults, pResult)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pMetaInfos)
	rmcResponseStream.WriteListResult(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodGetMetas, rmcResponseBody)

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
