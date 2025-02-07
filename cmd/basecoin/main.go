package main

import (
	"os"

	"github.com/tepleton/basecoin/cmd/basecoin/commands"
	"github.com/tepleton/tmlibs/cli"
)

func main() {
	rt := commands.RootCmd

	rt.AddCommand(
		commands.InitCmd,
		commands.StartCmd,
		commands.RelayCmd,
		commands.UnsafeResetAllCmd,
		commands.VersionCmd,
	)

	cmd := cli.PrepareMainCmd(rt, "BC", os.ExpandEnv("$HOME/.basecoin"))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
