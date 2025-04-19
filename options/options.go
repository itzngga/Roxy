package options

import (
	"database/sql"
	"errors"
	waTypes "go.mau.fi/whatsmeow/types"
	"time"
)

type LoginOptions int8

const (
	SCAN_QR   LoginOptions = 0
	PAIR_CODE LoginOptions = 1
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

	// OSInfo system name in client
	OSInfo string

	// LoginOptions constant of ScanQR or PairCode
	LoginOptions LoginOptions

	// HistorySync is used to synchronize message history
	HistorySync bool
	// AutoRejectCall allow to auto reject incoming calls
	AutoRejectCall bool

	// Bot General Settings

	// AllowFromPrivate allow messages from private
	AllowFromPrivate bool
	// AllowFromGroup allow message from groups
	AllowFromGroup bool
	// OnlyFromSelf allow only from self messages
	OnlyFromSelf bool
	// CommandSuggestion allow command suggestion
	CommandSuggestion bool
	// DebugMessage debug incoming message to console
	DebugMessage bool
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

func WithScanQRLogin() func(*Options) {
	return func(options *Options) {
		options.LoginOptions = SCAN_QR
	}
}

func WithPairCodeLogin() func(*Options) {
	return func(options *Options) {
		options.LoginOptions = PAIR_CODE
	}
}

func WithHistorySync() func(*Options) {
	return func(options *Options) {
		options.HistorySync = true
	}
}

func WithAutoRejectCall() func(*Options) {
	return func(options *Options) {
		options.AutoRejectCall = true
	}
}

func WithOSInfo(osInfo string) func(*Options) {
	return func(options *Options) {
		options.OSInfo = osInfo
	}
}

func WithDebugMessage() func(*Options) {
	return func(options *Options) {
		options.DebugMessage = true
	}
}

func NewDefaultOptions() *Options {
	return &Options{
		StoreMode:                   "sqlite",
		SqliteFile:                  "ROXY.DB",
		WithCommandLog:              true,
		AllowFromGroup:              true,
		AllowFromPrivate:            true,
		CommandSuggestion:           true,
		HistorySync:                 false,
		AutoRejectCall:              false,
		LoginOptions:                SCAN_QR,
		OSInfo:                      "Roxy",
		DebugMessage:                false,
		SendMessageTimeout:          time.Second * 30,
		CommandResponseCacheTimeout: time.Minute * 15,
	}
}

func (o *Options) Validate() error {
	if o.StoreMode != "postgres" && o.StoreMode != "sqlite" {
		return errors.New("error: invalid store mode")
	}
	if o.HostNumber != "" {
		if o.HostNumber[:2] == "08" {
			o.HostNumber = "628" + o.HostNumber[2:]
		}
		_, err := waTypes.ParseJID(o.HostNumber)
		if err != nil {
			return errors.New("error: invalid host number")
		}
	}

	if o.LogLevel == "" {
		o.LogLevel = "INFO"
	}

	if o.LogLevel != "INFO" && o.LogLevel != "ERROR" && o.LogLevel != "WARM" && o.LogLevel != "DEBUG" {
		return errors.New("error: invalid log level")
	}

	if o.WithSqlDB == nil && o.StoreMode == "postgres" && o.PostgresDsn == nil {
		return errors.New("error: postgresql dsn cannot be null")
	}

	if o.WithSqlDB == nil && o.StoreMode == "sqlite" && o.SqliteFile == "" {
		o.SqliteFile = "ROXY.DB"
	}

	if o.WithSqlDB == nil && o.SqliteFile == "" && o.PostgresDsn == nil {
		return errors.New("error: please specify sql.db or sqlite file or pg dsn")
	}

	if !o.AllowFromPrivate && !o.AllowFromGroup {
		return errors.New("error: please specify one of allow from private or group")
	}

	if o.PostgresDsn != nil {
		err := o.PostgresDsn.Validate()
		if err != nil {
			return err
		}
	}

	if o.LoginOptions == PAIR_CODE && o.HostNumber == "" {
		return errors.New("error: you must specify host number when using pair code login options")
	}

	if o.OSInfo == "" {
		o.OSInfo = "Roxy"
	}

	return nil
}
