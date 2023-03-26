package database

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/lib/pq"
)

func InsertMetaBinaryByDataStorePreparePostParamWithOwnerPID(dataStorePreparePostParam *nexproto.DataStorePreparePostParam, pid uint32) uint32 {
	var dataID uint32

	err := Postgres.QueryRow(`
		INSERT INTO mvdkts.meta_binaries (
			owner_pid,
			name,
			data_type,
			meta_binary,
			permission,
			del_permission,
			flag,
			period,
			tags,
			persistence_slot_id,
			extra_data
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING data_id`,
		pid,
		dataStorePreparePostParam.Name,
		dataStorePreparePostParam.DataType,
		dataStorePreparePostParam.MetaBinary,
		dataStorePreparePostParam.Permission.Permission,
		dataStorePreparePostParam.DelPermission.Permission,
		dataStorePreparePostParam.Flag,
		dataStorePreparePostParam.Period,
		pq.Array(dataStorePreparePostParam.Tags),
		dataStorePreparePostParam.PersistenceInitParam.PersistenceSlotId,
		pq.Array(dataStorePreparePostParam.ExtraData),
	).Scan(&dataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return dataID
}
