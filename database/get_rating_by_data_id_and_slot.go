package database

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func GetRatingByDataIDAndSlot(dataID uint32, slot int) (*datastore.DataStoreRatingInfoWithSlot, error) {
	rating := datastore.NewDataStoreRatingInfoWithSlot()
	rating.Slot = int8(slot)
	rating.Rating = datastore.NewDataStoreRatingInfo()

	err := Postgres.QueryRow(`
	SELECT
		total_value,
		count,
		initial_value
	FROM mvdkts.ratings WHERE data_id=$1 AND slot=$2`, dataID, slot).Scan(
		&rating.Rating.TotalValue,
		&rating.Rating.Count,
		&rating.Rating.InitialValue,
	)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	return rating, err
}
