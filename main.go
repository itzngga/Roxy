/*
Copyright © 2022 itzngga rangganak094@gmail.com. All rights reserved
*/
package main

import (
	"github.com/itzngga/goRoxy/internal"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	handler.NewDefaultMuxer()

}

func main() {
	base := internal.Base{}
	base.Init()
}
