package database

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func InsertRatingByDataIDAndDataStoreRatingInitParamWithSlot(dataID uint32, ratingInitParam *nexproto.DataStoreRatingInitParamWithSlot) error {
	_, err := Postgres.Exec(`
		INSERT INTO mvdkts.ratings (
			data_id,
			slot,
			flag,
			internal_flag,
			lock_type,
			initial_value,
			range_min,
			range_max,
			period_hour,
			period_duration,
			total_value,
			count
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 1)`,
		dataID,
		ratingInitParam.Slot,
		ratingInitParam.Param.Flag,
		ratingInitParam.Param.InternalFlag,
		ratingInitParam.Param.LockType,
		ratingInitParam.Param.InitialValue,
		ratingInitParam.Param.RangeMin,
		ratingInitParam.Param.RangeMax,
		ratingInitParam.Param.PeriodHour,
		ratingInitParam.Param.PeriodDuration,
		ratingInitParam.Param.InitialValue,
	)

	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return err
}
