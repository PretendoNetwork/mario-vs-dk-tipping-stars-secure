package main

import (
	"github.com/PretendoNetwork/mario-vs-dk-tipping-stars-secure/database"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
)

var logger = plogger.NewLogger()

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Warning("Error loading .env file")
	}

	database.ConnectAll()
}
