package options

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type PostgresDSN struct {
	Host     string
	Username string
	Password string
	Database string
	Port     string
	SslMode  string
	TimeZone string
}

func (dsn *PostgresDSN) Validate() error {
	if dsn.Host == "" {
		return errors.New("error: missing dsn host")
	}
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
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", dsn.Host, dsn.Port, dsn.Username, dsn.Password, dsn.Database, dsn.SslMode, dsn.TimeZone)
}

func NewPostgresDSN() *PostgresDSN {
	return &PostgresDSN{}
}

func (dsn *PostgresDSN) SetHost(host string) *PostgresDSN {
	dsn.Host = host
	return dsn
}

func (dsn *PostgresDSN) SetUsername(userName string) *PostgresDSN {
	dsn.Username = userName
	return dsn
}

func (dsn *PostgresDSN) SetPassword(password string) *PostgresDSN {
	dsn.Password = password
	return dsn
}

func (dsn *PostgresDSN) SetDatabase(database string) *PostgresDSN {
	dsn.Database = database
	return dsn
}

func (dsn *PostgresDSN) SetPort(port string) *PostgresDSN {
	dsn.Port = port
	return dsn
}

func (dsn *PostgresDSN) SetMode(sslMode string) *PostgresDSN {
	dsn.SslMode = sslMode
	return dsn
}

func (dsn *PostgresDSN) SetTimeZone(timeZone string) *PostgresDSN {
	dsn.TimeZone = timeZone
	return dsn
}

func (dsn *PostgresDSN) FromEnv() *PostgresDSN {
	dsn.Host = os.Getenv("PG_HOST")
	dsn.Port = os.Getenv("PG_PORT")
	dsn.Username = os.Getenv("PG_USERNAME")
	dsn.Password = os.Getenv("PG_PASSWORD")
	dsn.Database = os.Getenv("PG_DATABASE")
	dsn.SslMode = os.Getenv("PG_SSL_MODE")
	dsn.TimeZone = os.Getenv("PG_TIMEZONE")

	return dsn
}
