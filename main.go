package main

import (
	"github.com/alecthomas/kong"
	"github.com/n25a/portal/cli"
)

func main() {
	ctx := kong.Parse(&cli.CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
