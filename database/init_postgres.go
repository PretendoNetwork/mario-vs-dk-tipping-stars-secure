package database

import "github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/globals"

func initPostgres() {
	var err error

	_, err = Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS mvdkts`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres schema created")

	// TODO - Store ratings for meta binaries
	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS mvdkts.meta_binaries (
		data_id serial PRIMARY KEY,
		owner_pid integer,
		name text,
		data_type integer,
		meta_binary bytea,
		permission integer,
		del_permission integer,
		flag integer,
		period integer,
		tags text[],
		persistence_slot_id integer,
		extra_data text[]
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres tables created")
}
