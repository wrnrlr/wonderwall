package wonderwall

import ()

type Wall struct {
	ID      string        `json:"id"`
	Title   string        `json:"title"`
	Content []interface{} `json:"content"`
}

func (w *Wall) Key() Key {
	if w == nil {
		return Key("wall:")
	} else {
		return Key("wall:" + w.ID)
	}
}

type CreateWall interface {
	CreateWall(*Txn, *User) (*Wall, error)
}

type FindWallById interface {
	FindWallById(*Txn, string) (*Wall, error)
}

type UpdateWall interface {
	UpdateWall(*Txn, *Wall) error
}

type DeleteWall interface {
	DeleteWall(*Txn, *Wall) error
}

type Walls struct{}

func (w Walls) CreateWall(*Txn, *User) (*Wall, error)    { return nil, nil }
func (w Walls) FindWallById(*Txn, string) (*Wall, error) { return nil, nil }
func (w Walls) UpdateWall(*Txn, *Wall) error             { return nil }
func (w Walls) DeleteWall(*Txn, *Wall) error             { return nil }
