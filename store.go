package wonderwall

import (
	"bytes"
	"encoding/gob"
	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger/v2"
)

func Serialize(e interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(e)
	if err != nil {
		return []byte{}, err
	} else {
		return b.Bytes(), nil
	}
}

func Deserialize(v []byte, e interface{}) error {
	return gob.NewDecoder(bytes.NewReader(v)).Decode(e)
}

type Key []byte

func (k Key) Eq(other Key) bool {
	return bytes.Equal(k, other)
}

func (k *Key) Deserialize(b []byte) error {
	*k = b
	return nil
}

type Readable interface{ Deserialize([]byte) error }

type Writable interface {
	Serialize() ([]byte, error)
	Key() Key
}

type Index struct{ Primary, Secondary Key }

type Txn = badger.Txn

type Store struct {
	*badger.DB
	Index bleve.Index
}

type StoreConfig struct {
	Path       string
	InMemoryDB bool
}

func NewStore(conf StoreConfig) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(conf.Path).WithInMemory(conf.InMemoryDB))
	if err != nil {
		return nil, err
	}
	mappings := bleve.NewIndexMapping()
	WallMapping(mappings)
	index, err := bleve.NewMemOnly(mappings)
	if err != nil {
		return nil, err
	}
	return &Store{db, index}, nil
}

func (s *Store) Get(tx *Txn, k Key, r Readable) error {
	i, err := tx.Get(k)
	if err != nil {
		return err
	}
	return i.Value(r.Deserialize)
}

func (s *Store) Set(txn *Txn, o Writable) error {
	b, err := o.Serialize()
	if err != nil {
		return err
	}
	err = txn.Set(o.Key(), b)
	if err != nil {
		return err
	}
	return nil
}

var keyOnlyIterator = badger.IteratorOptions{PrefetchValues: false, PrefetchSize: 100, Reverse: false, AllVersions: false}

func (s *Store) Keys(txn *Txn, k Key, keys *[]Key) error {
	it := txn.NewIterator(keyOnlyIterator)
	defer it.Close()
	for it.Seek(k); it.ValidForPrefix(k); it.Next() {
		*keys = append(*keys, it.Item().Key())
	}
	return nil
}

func (s *Store) Delete(tx *Txn, k []byte) error {
	return tx.Delete(k)
}

// Stream values ...

func RebuildSearchIndex() {

}
