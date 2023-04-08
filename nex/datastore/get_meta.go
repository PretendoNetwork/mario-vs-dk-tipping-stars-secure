package nex_datastore

import (
	"fmt"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func GetMeta(err error, client *nex.Client, callID uint32, dataStoreGetMetaParam *datastore.DataStoreGetMetaParam) {
	var pMetaInfo *datastore.DataStoreMetaInfo
	var errorCode uint32

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	// TODO - Check dataStoreGetMetaParam.ResultOption and dataStoreGetMetaParam.AccessPassword?
	if dataStoreGetMetaParam.DataID == 0 {
		// * User related data?
		if dataStoreGetMetaParam.PersistenceTarget.PersistenceSlotID == 0 {
			pMetaInfo, errorCode = getMetaProfileInfo(dataStoreGetMetaParam)
		} else if dataStoreGetMetaParam.PersistenceTarget.PersistenceSlotID == 1 {
			pMetaInfo, errorCode = getMetaTipBucket(dataStoreGetMetaParam)
		} else {
			globals.Logger.Warning(fmt.Sprintf("UNKNOWN SLOT ID %d", dataStoreGetMetaParam.PersistenceTarget.PersistenceSlotID))
			// TODO - Send an error?
		}
	} else {
		globals.Logger.Warning(fmt.Sprintf("UNKNOWN TYPE %d", dataStoreGetMetaParam.DataID))
		// TODO - Send an error?
	}

	if errorCode == 0 {
		rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

		rmcResponseStream.WriteStructure(pMetaInfo)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodGetMeta, rmcResponseBody)
	} else {
		rmcResponse.SetError(errorCode)
	}

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

func getMetaProfileInfo(dataStoreGetMetaParam *datastore.DataStoreGetMetaParam) (*datastore.DataStoreMetaInfo, uint32) {
	metaBinary := database.GetMetaBinaryByTypeAndOwnerPIDAndSlotID(122, dataStoreGetMetaParam.PersistenceTarget.OwnerID, 0)
	pMetaInfo := datastore.NewDataStoreMetaInfo()

	if metaBinary.DataID == 0 {
		// * Meta binary doesn't exist
		return pMetaInfo, nex.Errors.DataStore.NotFound
	}

	pMetaInfo.DataID = uint64(metaBinary.DataID)
	pMetaInfo.OwnerID = metaBinary.OwnerPID
	pMetaInfo.Size = 0
	pMetaInfo.Name = metaBinary.Name
	pMetaInfo.DataType = metaBinary.DataType
	pMetaInfo.MetaBinary = metaBinary.Buffer
	pMetaInfo.Permission = datastore.NewDataStorePermission()
	pMetaInfo.Permission.Permission = 0
	pMetaInfo.Permission.RecipientIds = make([]uint32, 0)
	pMetaInfo.DelPermission = datastore.NewDataStorePermission()
	pMetaInfo.DelPermission.Permission = 3
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
	pMetaInfo.Ratings = make([]*datastore.DataStoreRatingInfoWithSlot, 0)

	return pMetaInfo, 0
}

func getMetaTipBucket(dataStoreGetMetaParam *datastore.DataStoreGetMetaParam) (*datastore.DataStoreMetaInfo, uint32) {
	metaBinary := database.GetMetaBinaryByTypeAndOwnerPIDAndSlotID(123, dataStoreGetMetaParam.PersistenceTarget.OwnerID, 1)
	pMetaInfo := datastore.NewDataStoreMetaInfo()

	if metaBinary.DataID == 0 {
		// * Meta binary doesn't exist
		return pMetaInfo, nex.Errors.DataStore.NotFound
	}

	// TODO - Check errors
	tipBucketExtraTipTotal, _ := database.GetRatingByDataIDAndSlot(metaBinary.DataID, 0)
	tipBucketExtraTipCount, _ := database.GetRatingByDataIDAndSlot(metaBinary.DataID, 1)
	tipBucketPlayCount, _ := database.GetRatingByDataIDAndSlot(metaBinary.DataID, 2)

	//tipBucketExtraTipTotal.Rating.Count = 20
	//tipBucketExtraTipTotal.Rating.TotalValue = 50 // stars
	//tipBucketExtraTipCount.Rating.Count = 30
	//tipBucketExtraTipCount.Rating.TotalValue = 100 // people
	//tipBucketPlayCount.Rating.Count = 40
	//tipBucketPlayCount.Rating.TotalValue = 150

	ratings := make([]*datastore.DataStoreRatingInfoWithSlot, 0, 3)

	ratings = append(ratings, tipBucketExtraTipTotal)
	ratings = append(ratings, tipBucketExtraTipCount)
	ratings = append(ratings, tipBucketPlayCount)

	pMetaInfo.DataID = uint64(metaBinary.DataID)
	pMetaInfo.OwnerID = metaBinary.OwnerPID
	pMetaInfo.Size = 0
	pMetaInfo.Name = metaBinary.Name
	pMetaInfo.DataType = metaBinary.DataType
	pMetaInfo.MetaBinary = metaBinary.Buffer
	pMetaInfo.Permission = datastore.NewDataStorePermission()
	pMetaInfo.Permission.Permission = 0
	pMetaInfo.Permission.RecipientIds = make([]uint32, 0)
	pMetaInfo.DelPermission = datastore.NewDataStorePermission()
	pMetaInfo.DelPermission.Permission = 3
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
	pMetaInfo.Ratings = ratings

	return pMetaInfo, 0
}
