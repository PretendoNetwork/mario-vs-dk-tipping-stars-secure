package nex_datastore

import (
	"fmt"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func GetMeta(err error, client *nex.Client, callID uint32, dataStoreGetMetaParam *nexproto.DataStoreGetMetaParam) {
	var pMetaInfo *nexproto.DataStoreMetaInfo
	var errorCode uint32

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)

	// TODO - Check dataStoreGetMetaParam.ResultOption and dataStoreGetMetaParam.AccessPassword?
	if dataStoreGetMetaParam.DataID == 0 {
		// * Mii data?
		pMetaInfo, errorCode = getMetaMiiData(dataStoreGetMetaParam)
	} else {
		globals.Logger.Warning(fmt.Sprintf("UNKNOWN TYPE %d", dataStoreGetMetaParam.DataID))
		// TODO - Send an error?
	}

	if errorCode == 0 {
		rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

		rmcResponseStream.WriteStructure(pMetaInfo)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(nexproto.DataStoreMethodGetMeta, rmcResponseBody)
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

func getMetaMiiData(dataStoreGetMetaParam *nexproto.DataStoreGetMetaParam) (*nexproto.DataStoreMetaInfo, uint32) {
	metaBinary := database.GetMetaBinaryByTypeAndOwnerPIDAndSlotID(122, dataStoreGetMetaParam.PersistenceTarget.OwnerID, uint8(dataStoreGetMetaParam.PersistenceTarget.PersistenceSlotID))
	pMetaInfo := nexproto.NewDataStoreMetaInfo()

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
	pMetaInfo.Permission = nexproto.NewDataStorePermission()
	pMetaInfo.Permission.Permission = 0
	pMetaInfo.Permission.RecipientIds = make([]uint32, 0)
	pMetaInfo.DelPermission = nexproto.NewDataStorePermission()
	pMetaInfo.DelPermission.Permission = 3
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

	return pMetaInfo, 0
}
