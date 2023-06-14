package options

import (
	"database/sql"
	"errors"
	"github.com/itzngga/Roxy/util"
	"time"
)

type Options struct {
	// HostNumber will use the first available device when null
	HostNumber string

	// StoreMode can be "postgres" or "sqlite"
	StoreMode string

	// LogLevel: "INFO", "ERROR", "WARN", "DEBUG"
	LogLevel string

	// This PostgresDsn Must add when StoreMode equal to "postgres"
	PostgresDsn PostgresDSN

	// This SqliteFile Generate "ROXY.DB" when it null
	SqliteFile string

	// WithSqlDB wrap with sql.DB interface
	WithSqlDB *sql.DB

	WithCommandLog              bool
	CommandResponseCacheTimeout time.Duration
	SendMessageTimeout          time.Duration
}

func NewDefaultOptions() Options {
	return Options{
		StoreMode:                   "sqlite",
		SqliteFile:                  "ROXY.DB",
		WithCommandLog:              true,
		SendMessageTimeout:          time.Second * 30,
		CommandResponseCacheTimeout: time.Minute * 15,
	}
}

func (o *Options) Validate() error {
	if !util.StringIsOnSlice(o.StoreMode, []string{"postgres", "sqlite"}) {
		return errors.New("error: invalid store mode")
	}
	if o.HostNumber != "" {
		_, ok := util.ParseJID(o.HostNumber)
		if !ok {
			return errors.New("error: invalid host number")
		}
	}
	if o.LogLevel == "" {
		o.LogLevel = "INFO"
	}
	if !util.StringIsOnSlice(o.LogLevel, []string{"INFO", "ERROR", "WARN", "DEBUG"}) {
		return errors.New("error: invalid log level")
	}

	var nilPgDsn PostgresDSN
	if o.WithSqlDB == nil && o.StoreMode == "postgres" && o.PostgresDsn == nilPgDsn {
		return errors.New("error: postgresql dsn cannot be null")
	}

	if o.WithSqlDB == nil && o.StoreMode == "sqlite" && o.SqliteFile == "" {
		o.SqliteFile = "GoRoxy.DB"
	}

	if o.WithSqlDB != nil && o.SqliteFile == "" && o.PostgresDsn == nilPgDsn {
		return errors.New("error: please specify sql.db or sqlite file or pg dsn")
	}

	return nil
}
