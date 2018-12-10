package main

import (
	"github.com/stefanm8/zyla/cli"
)

func main() {
	helper := cli.RunFromCLI()
	helper.Handle()
}
