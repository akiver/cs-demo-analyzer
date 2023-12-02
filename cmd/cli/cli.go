package main

import (
	"os"

	"github.com/akiver/cs-demo-analyzer/pkg/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
