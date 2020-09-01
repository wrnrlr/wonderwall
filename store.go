package main

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v2"
)

func serialize(e interface{}) ([]byte, error) {
	var b bytes.Buffer; enc := gob.NewEncoder(&b); err := enc.Encode(e)
	if err != nil { return []byte{}, err } else { return b.Bytes(), nil }}

func deserialize(v []byte, e interface{}) error {
	return gob.NewDecoder(bytes.NewReader(v)).Decode(e)}

type Key []byte
type Readable interface { Deserialize([]byte) error }
type Writable interface { Serialize() ([]byte,error); Key() Key }

type Index struct { Primary, Secondary Key }

type Store struct { *badger.DB }
type Txn = badger.Txn

func (s *Store) Get(tx *Txn, k Key, r Readable) error {
	i, err := tx.Get(k); if err != nil { return err }
	return i.Value(r.Deserialize)}

func (s *Store) Set(txn *Txn, o Writable) error {
	b, err := o.Serialize(); if err != nil { return err }
	if err := txn.Set(o.Key(), b); err != nil { return err }
	return nil }

var onlyKeysIterator = badger.IteratorOptions{PrefetchValues:false,PrefetchSize:100,Reverse:false,AllVersions:false,}
func (s *Store) Index(txn *Txn, k Key, keys *[]Key) error {
	it := txn.NewIterator(onlyKeysIterator); defer it.Close()
	for it.Seek(k); it.ValidForPrefix(k); it.Next() {
		*keys = append(*keys, it.Item().Key())}
	return nil}

func (s *Store) Delete(tx *Txn, k []byte) error { return tx.Delete(k) }