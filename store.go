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

type Store struct { *badger.DB }

func (s *Store) Get(k Key, r Readable) error {
	return s.View(func(tx *badger.Txn) error {
		if i, err := tx.Get(k); err != nil {
			return err
		} else {
			return i.Value(func(v []byte) error { return r.Deserialize(v) })
		}
	})
}

func (s *Store) Set(o Writable) error {
	b, err := o.Serialize(); if err != nil { return err }
	return s.DB.Update(func(txn *badger.Txn) error {
		if err := txn.Set(o.Key(), b); err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) Delete(k []byte) error {
	return s.Update(func(tx *badger.Txn) error { return tx.Delete(k) })
}

var onlyKeysIterator = badger.IteratorOptions{PrefetchValues:false,PrefetchSize:100,Reverse:false,AllVersions:false,}
func (s *Store) Index(k Key) (keys []Key, err error) {
	err = s.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(onlyKeysIterator)
		defer it.Close()
		for it.Seek(k); it.ValidForPrefix(k); it.Next() {
			keys = append(keys, it.Item().Key())
		}
		return nil
	})
	return keys, err
}
