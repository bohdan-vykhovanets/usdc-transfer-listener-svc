package main

import (
	"os"

	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
