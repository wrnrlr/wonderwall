package wonderwall

import "time"

type Session struct {
	ID        Token
	UserID    string
	CreatedAt time.Time
}

type CreateSession interface {
	CreateSession(*Txn, string) (*Session, error)
}

type Sessions struct{}

func (s Sessions) CreateSession(*Txn, string) (*Session, error) {
	return nil, nil
}
