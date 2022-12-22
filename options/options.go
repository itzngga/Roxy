package options

import (
	"github.com/itzngga/goRoxy/util"
	"time"
)

type Options struct {
	HostNumber string
	StoreMode  string
	LogLevel   string

	PostgresDsn string
	SqliteFile  string

	WithCommandCooldown bool
	WithCommandLog      bool
	WithBuiltIn         bool
	WithHelpCommand     bool

	HelpTitle       string
	HelpDescription string
	HelpFooter      string

	CommandCooldownTimeout      time.Duration
	CommandResponseCacheTimeout time.Duration
	SendMessageTimeout          time.Duration
}

func NewDefaultOptions() Options {
	return Options{
		StoreMode:                   "postgres",
		PostgresDsn:                 "user=akutansi password=root123 dbname=akutansi port=5432 sslmode=disable TimeZone=Asia/Jakarta",
		WithCommandCooldown:         false,
		WithCommandLog:              true,
		WithBuiltIn:                 true,
		WithHelpCommand:             true,
		SendMessageTimeout:          time.Second * 30,
		CommandResponseCacheTimeout: time.Minute * 15,

		HelpTitle:       "*GoRoxy BOT v1.2*",
		HelpDescription: "BOT ini dibuat dengan tujuan pembelajaran.\n\nHarap pilih salah satu",
		HelpFooter:      "@itzngga",
	}
}

func (o *Options) Validate() {
	if !util.StringIsOnSlice(o.StoreMode, []string{"postgres", "sqlite"}) {
		panic("error: invalid store mode")
	}
	if o.HostNumber != "" {
		_, ok := util.ParseJID(o.HostNumber)
		if !ok {
			panic("error: invalid host number")
		}
	}
	if o.LogLevel == "" {
		o.LogLevel = "INFO"
	}
	if !util.StringIsOnSlice(o.LogLevel, []string{"INFO", "ERROR", "WARN", "DEBUG"}) {
		panic("error: invalid log level")
	}
}
