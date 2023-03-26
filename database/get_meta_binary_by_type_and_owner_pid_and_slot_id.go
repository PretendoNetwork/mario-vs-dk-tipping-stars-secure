package database

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/types"
	"github.com/lib/pq"
)

func GetMetaBinaryByTypeAndOwnerPIDAndSlotID(dataType uint16, pid uint32, slotID uint8) *types.MetaBinary {
	metaBinary := types.NewMetaBinary()

	err := Postgres.QueryRow(`
	SELECT
	data_id,
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
	FROM mvdkts.meta_binaries WHERE data_type=$1 AND owner_pid=$2 AND persistence_slot_id=$3`, dataType, pid, slotID).Scan(
		&metaBinary.DataID,
		&metaBinary.OwnerPID,
		&metaBinary.Name,
		&metaBinary.DataType,
		&metaBinary.Buffer,
		&metaBinary.Permission,
		&metaBinary.DeletePermission,
		&metaBinary.Flag,
		&metaBinary.Period,
		pq.Array(&metaBinary.Tags),
		&metaBinary.PersistenceSlotID,
		pq.Array(&metaBinary.ExtraData),
	)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	return metaBinary
}
