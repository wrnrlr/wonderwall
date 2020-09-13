package main

type User struct {
	ID           string
	Email        Email
	PasswordHash PasswordHash
	Name         string
}

func (u *User) Key() Key {
	if u == nil {
		return Key("user:")
	} else {
		return Key("user:" + u.ID)
	}
}

type FindUserByEmail interface {
	FindUserByEmail(*Txn, Email) (*User, error)
}
type FindUserById interface {
	FindUserById(*Txn, string) (*User, error)
}
type CreateUser interface{ CreateUser(*Txn, *User) error }
type UpdateUser interface{ UpdateUser(*Txn, *User) error }
type DeleteUser interface{ DeleteUser(*Txn, *User) error }

type UserService interface {
	CreateUser
	UpdateUser
	DeleteUser
	FindUserById
	FindUserByEmail
}

type Users struct{ users []*Users }

func (s Users) CreateUser(*Txn, *User) error {
	return nil
}

func (s Users) UpdateUser(*Txn, *User) error               { return nil }
func (s Users) DeleteUser(*Txn, *User) error               { return nil }
func (s Users) FindUserById(*Txn, string) (*User, error)   { return nil, nil }
func (s Users) FindUserByEmail(*Txn, Email) (*User, error) { return nil, nil }
