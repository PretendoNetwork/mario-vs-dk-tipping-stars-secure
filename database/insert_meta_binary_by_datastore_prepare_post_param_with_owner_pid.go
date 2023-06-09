package database

import (
	"time"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/lib/pq"
)

func InsertMetaBinaryByDataStorePreparePostParamWithOwnerPID(dataStorePreparePostParam *datastore.DataStorePreparePostParam, pid uint32) uint32 {
	var dataID uint32

	now := time.Now().UnixNano()
	expireTime := time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC).UnixNano()

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
			extra_data,
			creation_time,
			updated_time,
			referred_time,
			expire_time
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING data_id`,
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
		now,
		now,
		now,
		expireTime,
	).Scan(&dataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return dataID
}
