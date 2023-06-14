# Roxy

Command Handler Abstraction for [whatsmeow](https://github.com/tulir/whatsmeow)

# Installation
```bash 
go get github.com/itzngga/Roxy
```
- You need ffmpeg binary for generate image/video thumbnail

# Get Started
```go
package main

import (
	_ "github.com/itzngga/Roxy/examples/cmd"
	"log"

	"github.com/itzngga/Roxy/core"
	"github.com/itzngga/Roxy/options"
	_ "github.com/mattn/go-sqlite3"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := core.NewGoRoxyBase(options.NewDefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}


```
# Config
#### default
```go
app := core.NewGoRoxyBase(options.NewDefaultOptions())
```
#### custom
```go
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
```
### PostgresSQL
```go
options := options.Options{
	StoreMode: "postgres"
	PostgresDsn: options.NewPostgresDSN("username", "password", "dbname", "port", "disable", "Asia/Jakarta")
}
app := core.NewGoRoxyBase(options)
```

### Sqlite
```go
options := options.Options{
	StoreMode: "sqlite"
	SqliteFile: "store.db"
}
app := core.NewGoRoxyBase(options)
```

# Add a Command
create a simple command with:
### command/hello_world.go
```go
package cmd

import (
	"fmt"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

var speed = &command.Command{
	Name:        "speed",
	Description: "Testing speed",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		t := time.Now()
		ctx.SendReplyMessage("wait...")
		return ctx.GenerateReplyMessage(fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}

func init() {
	embed.Commands.Add(speed)
}

```

# Documentation
[DOC](https://github.com/itzngga/Roxy/tree/master/DOC.md)
# Example
[Example](https://github.com/itzngga/Roxy/tree/master/examples)
# Helper/Util
[UTIL](https://github.com/itzngga/Roxy/tree/master/util)

# License
[GNU](https://github.com/itzngga/Roxy/blob/master/LICENSE)

# Bugs
Please submit an issue when Race Condition detected or anything

# Contribute
Pull Request are pleased to