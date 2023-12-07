<p align="center">
	<img src="https://c.tenor.com/wA8TRoy6bQoAAAAd/roxy-migurdia-mushoku-tensei.gif" width="200" height="180"/>
</p>

<p align="center">
	<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/ItzNgga/Roxy?style=flat-square">
	<img alt="GitHub forks" src="https://img.shields.io/github/forks/ItzNgga/Roxy?style=flat-square">
	<img alt="GitHub watchers" src="https://img.shields.io/github/watchers/ItzNgga/Roxy?style=flat-square">
	<img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/ItzNgga/Roxy?style=flat-square">
</p>

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

	"github.com/itzngga/Roxy"
	"github.com/itzngga/Roxy/options"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := roxy.NewRoxyBase(options.NewDefaultOptions())
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
app, err := roxy.NewRoxyBase(options.NewDefaultOptions())
if err != nil {
    log.Fatal(err)
}
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
```
### PostgresSQL
#### from env
```go
package main

import (
	roxy "github.com/itzngga/Roxy"
	_ "github.com/itzngga/Roxy/examples/cmd"
	"github.com/itzngga/Roxy/options"
	"github.com/joho/godotenv"

	"log"

	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	// Required ENV
	// PG_HOST : postgresql host
	// PG_PORT : postgresql port
	// PG_USERNAME : postgresql username
	// PG_PASSWORD : postgresql password
	// PG_DATABASE : postgresql database

	opt := options.NewDefaultOptions()
	opt.StoreMode = "postgres"
	opt.PostgresDsn = options.NewPostgresDSN().FromEnv()

	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
```
#### default parser
```go
package main

import (
	roxy "github.com/itzngga/Roxy"
	_ "github.com/itzngga/Roxy/examples/cmd"
	"github.com/itzngga/Roxy/options"

	"log"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	pg := options.NewPostgresDSN()
	pg.SetHost("localhost")
	pg.SetPort("4321")
	pg.SetUsername("postgres")
	pg.SetPassword("root123")
	pg.SetDatabase("roxy")

	opt := options.NewDefaultOptions()
	opt.StoreMode = "postgres"
	opt.PostgresDsn = pg

	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
```

### Sqlite
```go
package main

import (
	roxy "github.com/itzngga/Roxy"
	_ "github.com/itzngga/Roxy/examples/cmd"
	"github.com/itzngga/Roxy/options"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	opt := options.NewDefaultOptions()
	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}

```

# Add a Command
create a simple command with:
### command/hello_world.go
```go
package cmd

import (
	"fmt"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	"github.com/itzngga/Roxy"
	"github.com/itzngga/Roxy/context"
	"time"
)

func init() {
	roxy.Commands.Add(speed)
}

var speed = &roxy.Command{
	Name:        "speed",
	Description: "Testing speed",
	RunFunc: func(ctx *context.Ctx) *waProto.Message {
		t := time.Now()
		ctx.SendReplyMessage("wait...")
		return ctx.GenerateReplyMessage(fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}
```

# Create Question State
example with media question state
```go
package media

import (
	"github.com/itzngga/Leficious/src/cmd/constant"
	"github.com/itzngga/Roxy"
	"github.com/itzngga/Roxy/context"
	"github.com/itzngga/Roxy/util/cli"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	"log"
)

func init() {
	roxy.Commands.Add(ocr)
}

var ocr = &roxy.Command{
	Name:        "ocr",
	Category:    "media",
	Description: "Scan text on images",
	context: func(ctx *context.Ctx) *waProto.Message {
		var captured *waProto.Message
		ctx.NewUserQuestion().
			CaptureMediaQuestion("Please send/reply a media message", &captured).
			Exec()

		result, err := ctx.Download(captured, false)
		if err != nil {
			log.Fatal(err)
		}
		res := cli.ExecPipeline("tesseract", result, "stdin", "stdout", "-l", "ind", "--oem", "1", "--psm", "3", "-c", "preserve_interword_spaces=1")

		return ctx.GenerateReplyMessage(string(res))
	},
}
```
# Example
currently available example project in [Lara](https://github.com/itzngga/Lara)
# Documentation
[DOC](https://github.com/itzngga/Roxy/tree/master/DOC.md)
# License
[GNU](https://github.com/itzngga/Roxy/blob/master/LICENSE)

# Contribute
Pull Request are pleased to