package options

import (
	"errors"
	"github.com/itzngga/roxy/util"
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

	WithCommandLog              bool
	CommandResponseCacheTimeout time.Duration
	SendMessageTimeout          time.Duration
}

func NewDefaultOptions() Options {
	return Options{
		//HostNumber:                  "081297980063",
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
	if o.StoreMode == "postgres" && o.PostgresDsn == nilPgDsn {
		return errors.New("error: postgresql dsn cannot be null")
	}

	if o.StoreMode == "sqlite" && o.SqliteFile == "" {
		o.SqliteFile = "GoRoxy.DB"
	}

	return nil
}
