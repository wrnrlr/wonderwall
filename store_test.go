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
	var (keys []Key; dummy obj; o1Copy, o2Copy obj)
	o1 := obj{"hello"}; o2 := obj{"world"}; s := MemStore()
	assert.Nil(t, s.Update(func(txn *Txn)error{return s.Set(txn, &o1)}))
	assert.Nil(t, s.Update(func(txn *Txn)error{return s.Set(txn, &o2)}))
	assert.Nil(t, s.View(func(txn *Txn)error{return s.Index(txn, Key("obj:"), &keys)}))
	assert.Equal(t, 2, len(keys))
	assert.Nil(t, s.View(func(txn *Txn)error{return s.Get(txn, o1.Key(), &o1Copy)}))
	assert.Nil(t, s.View(func(txn *Txn)error{return s.Get(txn, o2.Key(), &o2Copy)}))
	assert.NotEqual(t, o1Copy, o2Copy)
	assert.Equal(t, o1, o1Copy)
	assert.NotNil(t, s.View(func(txn *Txn)error{return s.Get(txn, Key("unknown"), &dummy)}))
	assert.Nil(t, s.Update(func(txn *Txn)error{return s.Delete(txn, o2.Key())}))
	assert.NotNil(t, s.View(func(txn *Txn)error{return s.Get(txn, o2.Key(), &dummy)}))
	keys = nil;assert.Nil(t, s.View(func(txn *Txn)error{return s.Index(txn, Key("obj:"), &keys)}))
	assert.Equal(t, 1, len(keys))
}