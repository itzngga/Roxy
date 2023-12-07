package options

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type PostgresDSN struct {
	Host     string
	Username string
	Password string
	Database string
	Port     string
	SslMode  string
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
	} else {
		_, err := strconv.Atoi(dsn.Port)
		if err != nil {
			return errors.New("error: invalid given dsn port")
		}
	}
	if dsn.SslMode == "" {
		dsn.SslMode = "disable"
	}
	if dsn.SslMode != "disable" && dsn.SslMode != "verify-full" {
		return errors.New(fmt.Sprintf("error: invalid dsn ssl mode, given %s", dsn.SslMode))
	}

	return nil
}

func (dsn *PostgresDSN) GenerateDSN() string {
	return fmt.Sprintf("postgress://%s:%s@%s:%s/%s?sslmode=%s",
		dsn.Username,
		dsn.Password,
		dsn.Host,
		dsn.Port,
		dsn.Database,
		dsn.SslMode,
	)
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

func (dsn *PostgresDSN) FromEnv() *PostgresDSN {
	dsn.Host = os.Getenv("PG_HOST")
	dsn.Port = os.Getenv("PG_PORT")
	dsn.Username = os.Getenv("PG_USERNAME")
	dsn.Password = os.Getenv("PG_PASSWORD")
	dsn.Database = os.Getenv("PG_DATABASE")
	dsn.SslMode = os.Getenv("PG_SSL_MODE")

	return dsn
}
