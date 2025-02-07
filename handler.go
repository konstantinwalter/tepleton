package basecoin

import (
	"bytes"

	wrsp "github.com/tepleton/wrsp/types"
	crypto "github.com/tepleton/go-crypto"
	"github.com/tepleton/go-wire/data"

	"github.com/tepleton/basecoin/types"
)

// Handler is anything that processes a transaction
type Handler interface {
	CheckTx(ctx Context, store types.KVStore, tx Tx) (Result, error)
	DeliverTx(ctx Context, store types.KVStore, tx Tx) (Result, error)

	// TODO: flesh these out as well
	// SetOption(store types.KVStore, key, value string) (log string)
	// InitChain(store types.KVStore, vals []*wrsp.Validator)
	// BeginBlock(store types.KVStore, hash []byte, header *wrsp.Header)
	// EndBlock(store types.KVStore, height uint64) wrsp.ResponseEndBlock
}

// TODO: Context is a place-holder, soon we add some request data here from the
// higher-levels (like tell an app who signed).
// Trust me, we will need it like CallContext now...
type Context struct {
	sigs []crypto.PubKey
}

// TOTALLY insecure.  will redo later, but you get the point
func (c Context) AddSigners(keys ...crypto.PubKey) Context {
	return Context{
		sigs: append(c.sigs, keys...),
	}
}

func (c Context) GetSigners() []crypto.PubKey {
	return c.sigs
}

func (c Context) IsSignerAddr(addr []byte) bool {
	for _, pk := range c.sigs {
		if bytes.Equal(addr, pk.Address()) {
			return true
		}
	}
	return false
}

func (c Context) IsSignerKey(key crypto.PubKey) bool {
	for _, pk := range c.sigs {
		if key.Equals(pk) {
			return true
		}
	}
	return false
}

// Result captures any non-error wrsp result
// to make sure people use error for error cases
type Result struct {
	Data data.Bytes
	Log  string
}

func (r Result) ToWRSP() wrsp.Result {
	return wrsp.Result{
		Data: r.Data,
		Log:  r.Log,
	}
}
