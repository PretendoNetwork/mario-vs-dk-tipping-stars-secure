package database

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
)

func DeleteMetaBinaryByDataID(dataID uint32) error {
	_, err := Postgres.Exec(`DELETE FROM mvdkts.meta_binaries WHERE data_id=$1`, dataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return err
	}

	return nil
}
