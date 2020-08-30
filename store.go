package main

import (
	"github.com/dgraph-io/badger"
)

type Key []byte
type Readable interface { Deserialize([]byte) error }
type Writable interface { Serialize() []byte; Key() Key }

type Store struct { badger.DB }

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
	return s.DB.Update(func(txn *badger.Txn) error {
		if err := txn.Set(o.Key(), o.Serialize()); err != nil {
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
