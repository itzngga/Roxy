# goRoxy

a Golang version of Roxy WhatsApp Bot with Command Handler helper

# Installation

> go mod tidy

# Run
Normal run mode
> go run *.go

Run with race conditions' detector (DEBUG)
> go run --race *.go

# Environment
setup by copy the .env.example to .env

### PostgresSQL
`STORE_MODE=postgres`

### Sqlite
`STORE_MODE=sqlite`

`SQLITE_FILE=store.db`

### Command Cooldown Duration
`DEFAULT_COOLDOWN_SEC=5`

# Add a Command
create a simple command with:

### command/hello_world.go
```go
package command

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func HelloCommand() {
	AddCommand(&handler.Command{
		Name:        "hello",
		Aliases:     []string{"hai", "helo"},
		Description: "Command for Hello World!",
		Category:    handler.MiscCategory,
		RunFunc:     HelloRunFunc,
	})
}

func HelloRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	return util.SendReplyText(args.Evm, "Hello World!")
}
```
### command/zInit.go
```go
package command

import "github.com/itzngga/goRoxy/internal/handler"

var Commands []*handler.Command

func GenerateAllCommands() {
	HelloCommand() // <---- append new command here
}

func AddCommand(command *handler.Command) {
	Commands = append(Commands, command)
}
```

# Middlewares
middleware is function before RunFunc is executed

### Command middleware
is only this command middleware
```go
package command

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func HelloCommand() {
	AddCommand(&handler.Command{
		Name:        "hello",
		Aliases:     []string{"hai", "helo"},
		Description: "Command for Hello World!",
		Category:    handler.MiscCategory,
		RunFunc:     HelloRunFunc,
		Middleware:  HelloMiddleware,
	})
}

func HelloRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	return util.SendReplyText(args.Evm, "Hello World!")
}
func HelloMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	fmt.Println("Hi middleware!")
	return true
}
```
### Global middleware
all command run this middleware

### middleware/log.go
```go
package middleware

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func LogMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	fmt.Println("\n[CMD] Command : " + args.Cmd.Name)
	return true
}
```

### middleware/zInit.go
```go
package middleware

import "github.com/itzngga/goRoxy/internal/handler"

func GenerateAllMiddlewares() {
	AddMiddleware(LogMiddleware) // <-- append new middleware here
}

func AddMiddleware(mid handler.MiddlewareFunc) {
	handler.GlobalMiddleware = append(handler.GlobalMiddleware, mid)
}
```
# Helper/Util
[UTIL](https://github.com/itzngga/goRoxy/tree/master/util)

# License
[GNU](https://github.com/ItzNgga/goRoxy/blob/master/LICENSE)

# Bugs
Please submit an issue when Race Condition detected or anything

# Contribute
Pull Request are pleased to