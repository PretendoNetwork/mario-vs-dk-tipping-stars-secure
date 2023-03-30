package database

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
)

func UpdateRatingByDataIDAndSlot(dataID uint32, slot uint8, value int32) error {
	_, err := Postgres.Exec(`
		UPDATE mvdkts.ratings
		SET total_value=total_value+$1, count=count+1
		WHERE data_id=$2 AND slot=$3`,
		value,
		dataID,
		slot,
	)

	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return err
}
