package main

import (
	"mmplat/internal/commands"
	"os"
)

func main() {
	if err := commands.NewAppCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
