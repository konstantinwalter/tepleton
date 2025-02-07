package state

import (
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/types"
	eyes "github.com/tepleton/merkleeyes/client"
	"github.com/tepleton/tmlibs/log"
)

// CONTRACT: State should be quick to copy.
// See CacheWrap().
type State struct {
	chainID    string
	store      types.KVStore
	readCache  map[string][]byte // optional, for caching writes to store
	writeCache *types.KVCache    // optional, for caching writes w/o writing to store
	logger     log.Logger
}

func NewState(store types.KVStore) *State {
	return &State{
		chainID:    "",
		store:      store,
		readCache:  make(map[string][]byte),
		writeCache: nil,
		logger:     log.NewNopLogger(),
	}
}

func (s *State) SetLogger(l log.Logger) {
	s.logger = l
}

func (s *State) SetChainID(chainID string) {
	s.chainID = chainID
	s.store.Set([]byte("base/chain_id"), []byte(chainID))
}

func (s *State) GetChainID() string {
	if s.chainID != "" {
		return s.chainID
	}
	s.chainID = string(s.store.Get([]byte("base/chain_id")))
	return s.chainID
}

func (s *State) Get(key []byte) (value []byte) {
	if s.readCache != nil { //if not a cachewrap
		value, ok := s.readCache[string(key)]
		if ok {
			return value
		}
	}
	return s.store.Get(key)
}

func (s *State) Set(key []byte, value []byte) {
	if s.readCache != nil { //if not a cachewrap
		s.readCache[string(key)] = value
	}
	s.store.Set(key, value)
}

func (s *State) GetAccount(addr []byte) *types.Account {
	return types.GetAccount(s, addr)
}

func (s *State) SetAccount(addr []byte, acc *types.Account) {
	types.SetAccount(s, addr, acc)
}

func (s *State) CacheWrap() *State {
	cache := types.NewKVCache(s)
	return &State{
		chainID:    s.chainID,
		store:      cache,
		readCache:  nil,
		writeCache: cache,
		logger:     s.logger,
	}
}

// NOTE: errors if s is not from CacheWrap()
func (s *State) CacheSync() {
	s.writeCache.Sync()
}

func (s *State) Commit() wrsp.Result {
	switch s.store.(type) {
	case *eyes.Client:
		s.readCache = make(map[string][]byte)
		return s.store.(*eyes.Client).CommitSync()
	default:
		return wrsp.NewError(wrsp.CodeType_InternalError, "can only use Commit if store is merkleeyes")
	}

}
