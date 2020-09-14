package test

import (
	. "github.com/Almanax/wonderwall"
	"github.com/dgraph-io/badger/v2"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockStore() *Store {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic("failed to create test user")
	}
	return &Store{db}
}

func mockUser(name, email, pass string) *User {
	password, _ := Password(pass).HashPassword()
	return &User{xid.New(), Email(email), password, name}
}

func TestUsers(t *testing.T) {
	store := mockStore()
	users := Users{store}
	u := mockUser("alice", "alice@example.com", "Abcd1234!")
	createUser := func(u *User) error {
		return store.Update(func(txn *badger.Txn) error {
			return users.CreateUser(txn, u)
		})
	}
	assert.Nil(t, createUser(u))
	assert.Equal(t, DuplicateEmailErr, createUser(u))
	findUserById := func(id xid.ID) (u *User, err error) {
		store.View(func(txn *badger.Txn) error {
			u, err = users.FindUserById(txn, id)
			return err
		})
		return
	}
	_, err := findUserById(xid.New())
	assert.Equal(t, err, badger.ErrKeyNotFound)
	u1, err := findUserById(u.ID)
	assert.Nil(t, err)
	assert.True(t, u.Eq(u1))
	findUserByEmail := func(email Email) (u *User, err error) {
		store.View(func(txn *badger.Txn) error {
			u, err = users.FindUserByEmail(txn, email)
			return err
		})
		return
	}
	_, err = findUserByEmail("bob@example.com")
	assert.Equal(t, err, badger.ErrKeyNotFound)
	u2, err := findUserByEmail(u.Email)
	assert.Nil(t, err)
	assert.True(t, u.Eq(u2))
	deleteUser := func(u *User) error {
		return store.Update(func(txn *badger.Txn) error {
			return users.DeleteUser(txn, u)
		})
	}
	assert.Equal(t, deleteUser(mockUser("bob", "bob@example.com", "Hello1234!")), badger.ErrKeyNotFound)
	assert.Nil(t, deleteUser(u))
}
