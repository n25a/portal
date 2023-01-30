package main

import (
	"github.com/alecthomas/kong"
)

func main() {
	ctx := kong.Parse(&cli.CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
