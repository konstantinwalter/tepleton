package commands

import "github.com/tepleton/basecoin/plugins/ibc"

// returns a new IBC plugin to be registered with Basecoin
func NewIBCPlugin() *ibc.IBCPlugin {
	return ibc.New()
}
