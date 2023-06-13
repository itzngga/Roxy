package options

import (
	"errors"
	"fmt"
	"time"
)

type PostgresDSN struct {
	Username string
	Password string
	Database string
	Port     string
	SslMode  string
	TimeZone string
}

func (dsn *PostgresDSN) Validate() error {
	if dsn.Username == "" {
		return errors.New("error: missing dsn username")
	}
	if dsn.Password == "" {
		return errors.New("error: missing dsn password")
	}
	if dsn.Database == "" {
		return errors.New("error: missing dsn database")
	}
	if dsn.Port == "" {
		return errors.New("error: missing dsn port")
	}
	if dsn.SslMode == "" {
		dsn.SslMode = "disable"
	}
	if dsn.SslMode != "disable" && dsn.SslMode != "enable" {
		return errors.New(fmt.Sprintf("error: invalid dsn ssl mode, given %s", dsn.SslMode))
	}
	if dsn.TimeZone == "" {
		dsn.TimeZone = "Asia/Jakarta"
	} else {
		_, err := time.LoadLocation(dsn.TimeZone)
		if err != nil {
			return errors.New(fmt.Sprintf("error: invalid dsn timezone format, given %s", dsn.TimeZone))
		}
	}

	return nil
}

func (dsn *PostgresDSN) GenerateDSN() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dsn.Username, dsn.Password, dsn.Database, dsn.Port, dsn.SslMode, dsn.TimeZone)
}

func NewPostgresDSN(username string, password string, database string, port string, sslMode string, timeZone string) (PostgresDSN, error) {
	pDsn := PostgresDSN{
		Username: username,
		Password: password,
		Database: database,
		Port:     port,
		SslMode:  sslMode,
		TimeZone: timeZone,
	}

	err := pDsn.Validate()
	if err != nil {
		return pDsn, err
	}

	return pDsn, nil
}
