package config

import (
	"go.mau.fi/whatsmeow/store/sqlstore"
	"os"
)

func SqlStoreContainer() *sqlstore.Container {
	if os.Getenv("STORE_MODE") == "postgres" {
		container, err := sqlstore.New("postgres", CreatePostgresDsn(), NewWaDBLog())
		if err != nil {
			panic(err)
		}
		return container
	} else {
		container, err := sqlstore.New("sqlite3", CreateSqliteDsn(), NewWaDBLog())
		if err != nil {
			panic(err)
		}
		return container
	}
}
