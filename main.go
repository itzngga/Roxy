/*
Copyright © 2022 itzngga rangganak094@gmail.com. All rights reserved
*/
package main

import (
	"fmt"

	"github.com/itzngga/goRoxy/internal"
	"github.com/itzngga/goRoxy/internal/handler"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("xtconf")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
		panic("failed to read config file")
	}

	handler.NewDefaultMuxer()

}

func main() {
	base := internal.Base{}
	base.Init()
}
