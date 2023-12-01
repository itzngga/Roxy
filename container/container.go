package container

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-whatsapp/whatsmeow/store"
	"github.com/go-whatsapp/whatsmeow/store/sqlstore"
	waTypes "github.com/go-whatsapp/whatsmeow/types"
	waLog "github.com/go-whatsapp/whatsmeow/util/log"
	"github.com/itzngga/Roxy/options"
	"sync"
)

var containerPool sync.Pool

func init() {
	containerPool = sync.Pool{
		New: func() interface{} {
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
			container.clientJID = *device.JID
		}

		return device, err
	} else {
		devices, err := container.storeContainer.GetAllDevices()
		if err != nil {
			return nil, err
		}

		var device *store.Device
		for _, containerDevice := range devices {
			if containerDevice.JID.ToNonAD().User == hostNumber {
				device = containerDevice
				break
			}
		}

		if device == nil {
			container.NewDevice = true
			device = container.storeContainer.NewDevice()
		} else {
			container.clientJID = *device.JID
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
	} else if container.storeMode == "postgresql" {
		db, err := sql.Open("postgres", options.PostgresDsn.GenerateDSN())
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
		postgresql := sqlstore.NewWithDB(db, "postgres", waLog.Stdout("Database", "ERROR", true))
		err = postgresql.Upgrade()
		if err != nil {
			return nil, err
		}

		container.DB = db
		container.storeContainer = postgresql
		container.InitializeTables()

		return container, nil
	} else if container.storeMode == "sqlite" {
		db, err := sql.Open("sqlite3", options.SqliteFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
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
