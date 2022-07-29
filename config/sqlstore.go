package config

import "go.mau.fi/whatsmeow/store/sqlstore"

func SqlstoreContainer() *sqlstore.Container {
	container, err := sqlstore.New("postgres", CreatePostgresDsn(), NewWaDBLog())
	if err != nil {
		panic(err)
	}
	return container
}
