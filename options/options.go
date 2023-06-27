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
	PostgresDsn *PostgresDSN

	// This SqliteFile Generate "ROXY.DB" when it null
	SqliteFile string

	// WithSqlDB wrap with sql.DB interface
	WithSqlDB *sql.DB

	WithCommandLog              bool
	CommandResponseCacheTimeout time.Duration
	SendMessageTimeout          time.Duration

	// Bot General Settings

	// AllowFromPrivate allow messages from private
	AllowFromPrivate bool
	// AllowFromGroup allow message from groups
	AllowFromGroup bool
	// OnlyFromSelf allow only from self messages
	OnlyFromSelf bool
}

func New(options ...func(*Options)) (*Options, error) {
	option := &Options{}

	for _, op := range options {
		op(option)
	}

	err := option.Validate()
	if err != nil {
		return option, err
	}

	return option, nil
}

func WithHostNumber(hostNumber string) func(*Options) {
	return func(options *Options) {
		options.HostNumber = hostNumber
	}
}

func WithStoreMode(storeMode string) func(*Options) {
	return func(options *Options) {
		options.StoreMode = storeMode
	}
}

func WithLogLevel(logLevel string) func(*Options) {
	return func(options *Options) {
		options.LogLevel = logLevel
	}
}

func WithPostgresDSN(pgDsn *PostgresDSN) func(*Options) {
	return func(options *Options) {
		options.PostgresDsn = pgDsn
	}
}

func WithSqliteFile(sqliteFile string) func(*Options) {
	return func(options *Options) {
		options.SqliteFile = sqliteFile
	}
}

func WithSqlDB(sqlDB *sql.DB) func(*Options) {
	return func(options *Options) {
		options.WithSqlDB = sqlDB
	}
}

func WithCommandLog(cmdLog bool) func(*Options) {
	return func(options *Options) {
		options.WithCommandLog = cmdLog
	}
}

func WithCmdCacheTimeout(respCacheTimeout time.Duration) func(*Options) {
	return func(options *Options) {
		options.CommandResponseCacheTimeout = respCacheTimeout
	}
}

func WithSendMsgTimeout(sendMsgTimeout time.Duration) func(*Options) {
	return func(options *Options) {
		options.SendMessageTimeout = sendMsgTimeout
	}
}

func WithAllowFromPrivate(onlyFromPrivate bool) func(*Options) {
	return func(options *Options) {
		options.AllowFromPrivate = onlyFromPrivate
	}
}

func WithAllowFromGroup(onlyFromGroup bool) func(*Options) {
	return func(options *Options) {
		options.AllowFromGroup = onlyFromGroup
	}
}

func WithOnlyFromSelf(onlyFromSelf bool) func(*Options) {
	return func(options *Options) {
		options.OnlyFromSelf = onlyFromSelf
	}
}

func NewDefaultOptions() *Options {
	return &Options{
		StoreMode:                   "sqlite",
		SqliteFile:                  "ROXY.DB",
		WithCommandLog:              true,
		AllowFromGroup:              true,
		AllowFromPrivate:            true,
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

	if o.WithSqlDB == nil && o.StoreMode == "postgres" && o.PostgresDsn == nil {
		return errors.New("error: postgresql dsn cannot be null")
	}

	if o.WithSqlDB == nil && o.StoreMode == "sqlite" && o.SqliteFile == "" {
		o.SqliteFile = "GoRoxy.DB"
	}

	if o.WithSqlDB == nil && o.SqliteFile == "" && o.PostgresDsn == nil {
		return errors.New("error: please specify sql.db or sqlite file or pg dsn")
	}

	if !o.AllowFromPrivate && !o.AllowFromGroup {
		return errors.New("error: please specify one of allow from private or group")
	}

	return nil
}
