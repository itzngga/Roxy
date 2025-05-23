package cmd

import (
	"fmt"
	"time"

	roxy "github.com/itzngga/Roxy"
	"github.com/itzngga/Roxy/context"
)

func init() {
	speed := roxy.NewCommand("speed")
	speed.SetDescription("Testing latency")
	speed.UseCache(false)
	speed.SetRunFunc(speedFn)

	childNya := roxy.NewCommand("nya")
	childNya.SetDescription("Testing subcommand")
	childNya.SetRunFunc(nyaFn)

	speed.AddSubCommands(childNya)

	roxy.Commands.Add(speed)
}

func speedFn(ctx *context.Ctx) context.Result {
	nyow := time.Now()
	ctx.SendReplyMessage("Checking your connection speed~! (づ｡◕‿‿◕｡)づ")
	lawtency := time.Since(nyow).Milliseconds()
	return ctx.GenerateReplyMessage(fmt.Sprintf("Pong! Current connection latency is %d ms (づ｡◕‿‿◕｡)づ", lawtency))
}

func nyaFn(ctx *context.Ctx) context.Result {
	return ctx.GenerateReplyMessage("Nya!")
}
