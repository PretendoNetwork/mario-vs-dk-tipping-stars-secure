package database

import (
	"time"

	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/lib/pq"
)

func UpdateMetaBinaryByDataStoreChangeMetaParam(dataStoreChangeMetaParam *datastore.DataStoreChangeMetaParam) error {
	// TODO - Check dataStoreChangeMetaParam.ModifiesFlag
	// TODO - Check dataStoreChangeMetaParam.CompareParam

	now := time.Now().UnixNano()

	_, err := Postgres.Exec(`
		UPDATE mvdkts.meta_binaries
		SET
		name=$1,
		permission=$2,
		del_permission=$3,
		period=$4,
		meta_binary=$5,
		tags=$6,
		data_type=$7,
		updated_time=$8
		WHERE data_id=$9`,
		dataStoreChangeMetaParam.Name,
		dataStoreChangeMetaParam.Permission.Permission,
		dataStoreChangeMetaParam.DelPermission.Permission,
		dataStoreChangeMetaParam.Period,
		dataStoreChangeMetaParam.MetaBinary,
		pq.Array(dataStoreChangeMetaParam.Tags),
		dataStoreChangeMetaParam.DataType,
		now,
		dataStoreChangeMetaParam.DataID,
	)

	if err != nil {
		return err
	}

	return nil
}
