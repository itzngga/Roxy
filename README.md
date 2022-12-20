# goRoxy

a Golang version of Roxy WhatsApp Bot with Command Handler helper

# Installation
```bash 
go get github.com/itzngga/goRoxy
```

# Get Started
```go
package main

import (
	_ "github.com/itzngga/goRoxy/examples/cmd"

	"github.com/itzngga/goRoxy/core"
	"github.com/itzngga/goRoxy/options"
	_ "github.com/lib/pq"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := core.NewGoRoxyBase(options.NewDefaultOptions())

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
    HostNumber string
    StoreMode  string
    LogLevel   string
    
    PostgresDsn string
    SqliteFile  string
    
    WithCommandCooldown bool
    WithCommandLog      bool
    
    HelpTitle       string
    HelpDescription string
    HelpFooter      string
    
    CommandCooldownTimeout      time.Duration
    CommandResponseCacheTimeout time.Duration
    SendMessageTimeout          time.Duration
}
```
### PostgresSQL
```go
options := options.Options{
	StoreMode: "postgres"
	PostgresDsn: "user=goroxy password=test123 dbname=goroxy port=5432 sslmode=disable TimeZone=Asia/Jakarta"
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
	"github.com/itzngga/goRoxy/basic/categories"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

var speed = &command.Command{
	Name:        "speed",
	Aliases:     []string{"sp", "s"},
	Description: "Testing speed",
	Category:    categories.CommonCategory,
	RunFunc: func(c *whatsmeow.Client, params *command.RunFuncParams) *waProto.Message {
		t := time.Now()
		util.SendReplyMessage(c, params.Event, "ok, waitt...")
		return util.SendReplyText(params.Event, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}

func init() {
	embed.Commands.Add(speed)
}
```

# Documentation
[DOC](https://github.com/itzngga/goRoxy/tree/master/DOC.md)
# Example
[Example](https://github.com/itzngga/goRoxy/tree/master/examples)
# Helper/Util
[UTIL](https://github.com/itzngga/goRoxy/tree/master/util)

# License
[GNU](https://github.com/ItzNgga/goRoxy/blob/master/LICENSE)

# Bugs
Please submit an issue when Race Condition detected or anything

# Contribute
Pull Request are pleased to