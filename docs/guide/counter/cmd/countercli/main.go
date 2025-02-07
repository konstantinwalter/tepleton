package main

import (
	"os"

	"github.com/spf13/cobra"

	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/light-client/commands"
	"github.com/tepleton/light-client/commands/proofs"
	"github.com/tepleton/light-client/commands/proxy"
	"github.com/tepleton/light-client/commands/seeds"
	"github.com/tepleton/light-client/commands/txs"
	"github.com/tepleton/tmlibs/cli"

	bcmd "github.com/tepleton/basecoin/cmd/basecli/commands"
	bcount "github.com/tepleton/basecoin/docs/guide/counter/cmd/countercli/commands"
)

// BaseCli represents the base command when called without any subcommands
var BaseCli = &cobra.Command{
	Use:   "countercli",
	Short: "Light client for tepleton",
	Long: `Basecli is an version of tmcli including custom logic to
present a nice (not raw hex) interface to the basecoin blockchain structure.

This is a useful tool, but also serves to demonstrate how one can configure
tmcli to work for any custom wrsp app.
`,
}

func main() {
	commands.AddBasicFlags(BaseCli)

	// Prepare queries
	proofs.RootCmd.AddCommand(
		// These are default parsers, optional in your app
		proofs.TxCmd,
		proofs.KeyCmd,
		bcmd.AccountQueryCmd,

		// XXX IMPORTANT: here is how you add custom query commands in your app
		bcount.CounterQueryCmd,
	)

	// Prepare transactions
	proofs.TxPresenters.Register("base", bcmd.BaseTxPresenter{})
	txs.RootCmd.AddCommand(
		// This is the default transaction, optional in your app
		bcmd.SendTxCmd,

		// XXX IMPORTANT: here is how you add custom tx construction for your app
		bcount.CounterTxCmd,
	)

	// Set up the various commands to use
	BaseCli.AddCommand(
		commands.InitCmd,
		commands.ResetCmd,
		keycmd.RootCmd,
		seeds.RootCmd,
		proofs.RootCmd,
		txs.RootCmd,
		proxy.RootCmd,
	)

	cmd := cli.PrepareMainCmd(BaseCli, "CTL", os.ExpandEnv("$HOME/.countercli"))
	cmd.Execute()
}
