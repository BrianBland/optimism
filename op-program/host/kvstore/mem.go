package kvstore

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"slices"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// MemKV implements the KV store interface in memory, backed by a regular Go map.
// This should only be used in testing, as large programs may require more pre-image data than available memory.
// MemKV is safe for concurrent use.
type MemKV struct {
	sync.RWMutex
	m map[common.Hash][]byte
}

var _ KV = (*MemKV)(nil)

func NewMemKV() *MemKV {
	return &MemKV{m: make(map[common.Hash][]byte)}
}

type fixture struct {
	WitnessData map[common.Hash]string `json:"witnessData"`
}

func FromFixture(path string) (*MemKV, error) {
	var f fixture
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&f)
	memKV := NewMemKV()

	for k, v := range f.WitnessData {
		val, err := hex.DecodeString(v[2:])
		if err != nil {
			return nil, err
		}
		memKV.Put(k, val)
	}

	return memKV, nil
}

func (m *MemKV) Put(k common.Hash, v []byte) error {
	m.Lock()
	defer m.Unlock()
	m.m[k] = slices.Clone(v)
	return nil
}

func (m *MemKV) Get(k common.Hash) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.m[k]
	if !ok {
		return nil, ErrNotFound
	}
	return slices.Clone(v), nil
}
