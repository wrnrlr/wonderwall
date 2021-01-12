package wonderwall

import (
	"bytes"
	"github.com/dgraph-io/badger/v2"
	"github.com/rs/xid"
	"strings"
)

type User struct {
	ID           xid.ID
	Email        Email
	PasswordHash PasswordHash
	Name         string
	Roles        []Role
}

func (u *User) Eq(other *User) bool {
	return u.ID == other.ID && u.Email == other.Email &&
		bytes.Compare(u.PasswordHash, other.PasswordHash) == 0 && u.Name == other.Name
}

func (u *User) Key() Key                   { return userPrimaryIndex(u.ID) }
func (u *User) Serialize() ([]byte, error) { return Serialize(u) }
func (u *User) Deserialize(b []byte) error { return Deserialize(b, u) }
func (u *User) EmailKey() Key              { return userEmailIndex(u.Email) }

func userPrimaryIndex(id xid.ID) Key { return append([]byte("user:"), id.Bytes()...) }
func userEmailIndex(email Email) Key { return []byte("user:email:" + strings.ToLower(email.String())) }

type FindUserByEmail interface {
	FindUserByEmail(*Txn, Email) (*User, error)
}
type FindUserById interface {
	FindUserById(*Txn, xid.ID) (*User, error)
}
type CreateUser interface{ CreateUser(*Txn, *User) error }
type DeleteUser interface{ DeleteUser(*Txn, *User) error }

type UserService interface {
	CreateUser
	DeleteUser
	FindUserById
	FindUserByEmail
}

type Users struct {
	DB *Store
}

func (s Users) CreateUser(txn *Txn, u *User) error {
	var primaryKey Key
	err := s.DB.Get(txn, u.EmailKey(), &primaryKey)
	if err != badger.ErrKeyNotFound {
		return DuplicateEmailErr
	}
	if err := txn.Set(u.EmailKey(), u.Key()); err != nil {
		return err
	}
	return s.DB.Set(txn, u)
}

func (s Users) DeleteUser(txn *Txn, u *User) error {
	err := txn.Delete(u.EmailKey())
	if err != nil {
		return err
	}
	return s.DB.Delete(txn, u.Key())
}

func (s Users) FindUserById(txn *Txn, id xid.ID) (*User, error) {
	var u User
	key := userPrimaryIndex(id)
	err := s.DB.Get(txn, key, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s Users) FindUserByEmail(txn *Txn, email Email) (*User, error) {
	var u User
	var primaryKey Key
	err := s.DB.Get(txn, userEmailIndex(email), &primaryKey)
	if err != nil {
		return nil, err
	}
	err = s.DB.Get(txn, primaryKey, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
