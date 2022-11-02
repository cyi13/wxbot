package main

import (
	"embed"
	"wxbot/cmd"
)

//go:embed dll
var DLLFS embed.FS

func main() {
	cmd.Execute()
}
