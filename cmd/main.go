package main

import (
	"mmplat/internal/command"
	"os"
)

func main() {
	if err := command.NewAppCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
