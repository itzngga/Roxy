# goRoxy

a Golang version of Roxy WhastApp Bot with Command Handler helper

# Installation

> go mod tidy

# Run
Normal run mode
> go run *.go

Run with race conditions detector (DEBUG)
> go run --race *.go

# Environment
setup by copy the .env.example to .env

### PostgreSQL
`STORE_MODE=postgres`

### Sqlite
`STORE_MODE=sqlite`

`SQLITE_FILE=store.db`

### Command Cooldown Duration
`DEFAULT_COOLDOWN_SEC=5`

# License
[GNU](https://github.com/ItzNgga/goRoxy/blob/master/LICENSE)

# Bugs
Please submit a issue when Race Condition detected or anything

# Contribute
Pull Request are pleased to