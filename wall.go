package wonderwall

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"
	"github.com/rs/xid"
)

var (
	ErrWallTitleEmpty = errors.New("wall title empty")
)

type Wall struct {
	ID         xid.ID        `json:"id"`
	Title      string        `json:"title"`
	Owner      Key           `json:"owner"`
	CreatedAt  time.Time     `json:"createdAt"`
	ModifiedAt time.Time     `json:"modifiedAt"`
	AllowList  []Key         `json:"allowList"`
	Content    []interface{} `json:"content"`
}

func (w Wall) Validate() error {
	if len(w.Title) == 0 {
		return ErrWallTitleEmpty
	}
	return nil
}

func (w *Wall) Key() Key                   { return wallPrimaryIndex(w.ID) }
func (w *Wall) Serialize() ([]byte, error) { return Serialize(w) }
func (w *Wall) Deserialize(b []byte) error { return Deserialize(b, w) }

func wallPrimaryIndex(id xid.ID) Key { return append([]byte("wall:"), id.Bytes()...) }

type CreateWall interface {
	CreateWall(*Txn, string, Key) (*Wall, error)
}

type FindWallById interface {
	FindWallById(*Txn, xid.ID) (*Wall, error)
}

type UpdateWall interface{ UpdateWall(*Txn, *Wall) error }

type DeleteWall interface{ DeleteWall(*Txn, *Wall) error }

type Walls struct {
	DB *Store
}

func (s Walls) CreateWall(txn *Txn, title string, owner Key) (*Wall, error) {
	now := time.Now()
	w := Wall{ID: xid.New(), Title: title, Owner: owner, CreatedAt: now, ModifiedAt: now}
	err := s.DB.Set(txn, &w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (s Walls) FindWallById(txn *Txn, id xid.ID) (*Wall, error) {
	var u Wall
	key := userPrimaryIndex(id)
	err := s.DB.Get(txn, key, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s Walls) UpdateWall(txn *Txn, w *Wall) error {
	return s.DB.Set(txn, w)
}

func (s Walls) DeleteWall(txn *Txn, w *Wall) error {
	return s.DB.Delete(txn, w.Key())
}

func (s Walls) SearchWalls() error {
	return nil
}

func WallMapping(index *mapping.IndexMappingImpl) {
	mapping := bleve.NewDocumentMapping()
	title := bleve.NewTextFieldMapping()
	title.Analyzer = simple.Name
	mapping.AddFieldMappingsAt("title", title)
	owner := bleve.NewTextFieldMapping()
	owner.Analyzer = keyword.Name
	mapping.AddFieldMappingsAt("owner", owner)
	createdAt := bleve.NewDateTimeFieldMapping()
	mapping.AddFieldMappingsAt("createdAt", createdAt)
	modifiedAt := bleve.NewDateTimeFieldMapping()
	mapping.AddFieldMappingsAt("modifiedAt", modifiedAt)
	index.AddDocumentMapping("wall", mapping)
}

func GetWallHandler(db *Store, security Auth, walls *Walls) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query query.Query
		query = bleve.NewMatchAllQuery()
		search := bleve.NewSearchRequest(query)
		search.Fields = []string{"id", "title", "owner", "createdAt", "modifiedAt"}
		results, err := db.Index.Search(search)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		writeTmpl(w, "walls", results)
	}
}

func PostWallHandler(db *Store, security Auth, walls *Walls) http.HandlerFunc {
	type wallForm struct {
		Title string `json:"title"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := security.Principal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if p == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var f wallForm
		err = json.NewDecoder(r.Body).Decode(&f)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var wall *Wall
		err = db.Update(func(txn *Txn) error {
			wall, err = walls.CreateWall(txn, f.Title, p.Key())
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func PatchWallHandler(security Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteWallHandler(db *Store, security Auth, walls *Walls) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := security.Principal(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if p == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := parseID(r.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.Update(func(txn *Txn) error {
			w, err := walls.FindWallById(txn, id)
			if err != nil {
				return err
			}
			if !w.Owner.Eq(p.Key()) {
				return AuthorizationError
			}
			return walls.DeleteWall(txn, w)
		})
		if err == AuthorizationError {
			w.WriteHeader(http.StatusForbidden)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func parseID(url *url.URL) (xid.ID, error) {
	parts := strings.Split(url.Path, "/")
	return xid.FromString(parts[len(parts)-1])
}
