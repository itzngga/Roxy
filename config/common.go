package config

import (
	"fmt"

	"github.com/spf13/viper"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func NewWaDBLog() waLog.Logger {
	return waLog.Stdout("Database", "INFO", true)
}
func NewWaLog() waLog.Logger {
	return waLog.Stdout("Main", "ERROR", true)
}

func CreatePostgresDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		viper.GetString("XT_DB_HOST"),
		viper.GetString("XT_DB_USER"),
		viper.GetString("XT_DB_PASSWORD"),
		viper.GetString("XT_DB_NAME"),
		viper.GetString("XT_DB_PORT"),
	)
}
