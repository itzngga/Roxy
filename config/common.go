package config

import (
	"fmt"
	"os"

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
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}

func CreateSqliteDsn() string {
	return fmt.Sprintf("file:%s?_foreign_keys=on", os.Getenv("SQLITE_FILE"))
}
