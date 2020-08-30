package main

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func MemStore() *Store {
	db, _ := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	return &Store{db}
}

type obj struct { ID string }
func (o *obj) Key() Key { if o == nil { return Key("obj:")} else { return Key("obj:"+o.ID) } }
func (o *obj) Serialize() ([]byte, error) { return serialize(o) }
func (o *obj) Deserialize(b []byte) error { return deserialize(b, o) }

func TestStore(t *testing.T) {
	o1 := obj{"hello"}; var o1Copy obj; s := MemStore()
	o2 := obj{"world"}; var o2Copy, dummy obj;
	assert.Nil(t, s.Set(&o1))
	assert.Nil(t, s.Set(&o2))
	keys, err := s.Index(Key("obj:"))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keys))
	assert.Nil(t, s.Get(o1.Key(), &o1Copy))
	assert.Nil(t, s.Get(o2.Key(), &o2Copy))
	assert.NotEqual(t, o1Copy, o2Copy)
	assert.NotNil(t, s.Get(Key("unknown"), &dummy))
	assert.Equal(t, o1, o1Copy)
	assert.Nil(t, s.Delete(o2.Key()))
	assert.NotNil(t, s.Get(o2.Key(), &dummy))
	keys, err = s.Index(Key("obj:"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
}