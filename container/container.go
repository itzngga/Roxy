package container

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/itzngga/Roxy/options"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var containerPool sync.Pool

func init() {
	containerPool = sync.Pool{
		New: func() any {
			return &Container{}
		},
	}
}

func AcquireContainer() *Container {
	return containerPool.Get().(*Container)
}

func ReleaseContainer(c *Container) {
	containerPool.Put(c)
}

type Container struct {
	clientJID      waTypes.JID
	DB             *sql.DB
	storeMode      string
	storeContainer *sqlstore.Container
	NewDevice      bool
}

func (container *Container) SetClientJID(jid waTypes.JID) {
	container.clientJID = jid
}

func (container *Container) ClientJID() waTypes.JID {
	return container.clientJID
}

func (container *Container) AcquireDevice(hostNumber string) (*store.Device, error) {
	if hostNumber == "" {
		device, err := container.storeContainer.GetFirstDevice()
		if device == nil {
			container.NewDevice = true
			device = container.storeContainer.NewDevice()
		} else {
			container.clientJID = *device.ID
		}

		return device, err
	} else {
		devices, err := container.storeContainer.GetAllDevices()
		if err != nil {
			return nil, err
		}

		var device *store.Device
		for _, containerDevice := range devices {
			if containerDevice.ID.ToNonAD().User == hostNumber {
				device = containerDevice
				break
			}
		}

		if device == nil {
			container.NewDevice = true
			device = container.storeContainer.NewDevice()
		} else {
			container.clientJID = *device.ID
		}

		return device, nil
	}
}

func NewContainer(options *options.Options) (*Container, error) {
	container := AcquireContainer()
	container.storeMode = options.StoreMode

	if options.WithSqlDB != nil {
		sqlDB := sqlstore.NewWithDB(options.WithSqlDB, options.StoreMode, waLog.Stdout("Database", "ERROR", true))
		err := sqlDB.Upgrade()
		if err != nil {
			return nil, err
		}

		container.DB = options.WithSqlDB
		container.storeContainer = sqlDB

		return container, nil
	} else if container.storeMode == "postgres" {
		db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(options.PostgresDsn.GenerateDSN())))
		postgresql := sqlstore.NewWithDB(db, "postgres", waLog.Stdout("Database", "ERROR", true))
		err := postgresql.Upgrade()
		if err != nil {
			return nil, err
		}

		container.DB = db
		container.storeContainer = postgresql
		container.InitializeTables()

		return container, nil
	} else if container.storeMode == "sqlite" {
		db, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?_foreign_keys=on&cache=shared", options.SqliteFile))
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		_, err = db.Exec("PRAGMA journal_mode=WAL;PRAGMA foreign_keys = ON;")
		if err != nil {
			return nil, fmt.Errorf("error sqlite: %w", err)
		}

		sqlite := sqlstore.NewWithDB(db, "sqlite3", waLog.Stdout("Database", "ERROR", true))
		err = sqlite.Upgrade()
		if err != nil {
			return nil, err
		}

		container.DB = db
		container.storeContainer = sqlite
		container.InitializeTables()

		return container, nil
	} else {
		ReleaseContainer(container)
		return nil, errors.New("error: invalid store mode")
	}
}
