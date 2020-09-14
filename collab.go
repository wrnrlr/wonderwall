package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
)

type CollabConfig struct {
	Debug     bool
	Workers   int
	Queue     int
	IOTimeout time.Duration
}

func nameConn(conn net.Conn) string {
	return conn.LocalAddr().String() + " > " + conn.RemoteAddr().String()
}

func WallCollab(conf CollabConfig, db *Store, walls interface {
	FindWallById
	UpdateWall
}) http.HandlerFunc {
	// Initialize netpoll instance. We will use it to be noticed about incoming
	// events from listener of user connections.
	poller, err := netpoll.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	var (
		// Make pool of X size, Y sized work queue and one pre-spawned
		// goroutine.
		pool = NewPool(conf.Workers, conf.Queue, 1)
		chat = NewChat(pool)
	)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Wall ws")
		// Check users...
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("can't upgrade http connection to websocket: %v", err)
			return
		}
		// NOTE: we wrap conn here to show that ws could work with any kind of
		// io.ReadWriter.
		safeConn := deadliner{conn, conf.IOTimeout}
		user := chat.Register(safeConn)
		// Create netpoll event descriptor for conn.
		// We want to handle only read events of it.
		desc := netpoll.Must(netpoll.HandleRead(conn))
		// Subscribe to events about conn.
		poller.Start(desc, func(ev netpoll.Event) {
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				// When ReadHup or Hup received, this mean that client has
				// closed at least write end of the connection or connections
				// itself. So we want to stop receive events about such conn
				// and remove it from the chat registry.
				poller.Stop(desc)
				chat.Remove(user)
				return
			}
			// Here we can read some new message from connection.
			// We can not read it right here in callback, because then we will
			// block the poller's inner loop.
			// We do not want to spawn a new goroutine to read single message.
			// But we want to reuse previously spawned goroutine.
			pool.Schedule(func() {
				if err := user.Receive(); err != nil {
					// When receive failed, we can only disconnect broken
					// connection and stop to receive events about it.
					poller.Stop(desc)
					chat.Remove(user)
				}
			})
		})
	}
}

type Object map[string]interface{}

type Request struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params Object `json:"params"`
}

type Response struct {
	ID     int    `json:"id"`
	Result Object `json:"result"`
}

type Error struct {
	ID    int    `json:"id"`
	Error Object `json:"error"`
}

type Collaborator struct {
	io   sync.Mutex
	conn io.ReadWriteCloser

	id   uint
	name string
	chat *Collective
}

// Receive reads next message from user's underlying connection.
// It blocks until full message received.
func (u *Collaborator) Receive() error {
	req, err := u.readRequest()
	if err != nil {
		u.conn.Close()
		return err
	}
	if req == nil {
		// Handled some control message.
		return nil
	}
	switch req.Method {
	case "rename":
		name, ok := req.Params["name"].(string)
		if !ok {
			return u.writeErrorTo(req, Object{
				"error": "bad params",
			})
		}
		prev, ok := u.chat.Rename(u, name)
		if !ok {
			return u.writeErrorTo(req, Object{
				"error": "already exists",
			})
		}
		u.chat.Broadcast("rename", Object{
			"prev": prev,
			"name": name,
			"time": timestamp(),
		})
		return u.writeResultTo(req, nil)
	case "publish":
		req.Params["author"] = u.name
		req.Params["time"] = timestamp()
		u.chat.Broadcast("publish", req.Params)
	default:
		return u.writeErrorTo(req, Object{
			"error": "not implemented",
		})
	}
	return nil
}

// readRequests reads json-rpc request from connection.
// It takes io mutex.
func (u *Collaborator) readRequest() (*Request, error) {
	u.io.Lock()
	defer u.io.Unlock()

	h, r, err := wsutil.NextReader(u.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(u.conn, ws.StateServerSide)(h, r)
	}

	req := &Request{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (u *Collaborator) writeErrorTo(req *Request, err Object) error {
	return u.write(Error{
		ID:    req.ID,
		Error: err,
	})
}

func (u *Collaborator) writeResultTo(req *Request, result Object) error {
	return u.write(Response{
		ID:     req.ID,
		Result: result,
	})
}

func (u *Collaborator) writeNotice(method string, params Object) error {
	return u.write(Request{
		Method: method,
		Params: params,
	})
}

func (u *Collaborator) write(x interface{}) error {
	w := wsutil.NewWriter(u.conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	u.io.Lock()
	defer u.io.Unlock()

	if err := encoder.Encode(x); err != nil {
		return err
	}

	return w.Flush()
}

func (u *Collaborator) writeRaw(p []byte) error {
	u.io.Lock()
	defer u.io.Unlock()

	_, err := u.conn.Write(p)

	return err
}

// Chat contains logic of user interaction.
type Collective struct {
	mu  sync.RWMutex
	seq uint
	us  []*Collaborator
	ns  map[string]*Collaborator

	pool *Pool
	out  chan []byte
}

func NewChat(pool *Pool) *Collective {
	chat := &Collective{
		pool: pool,
		ns:   make(map[string]*Collaborator),
		out:  make(chan []byte, 1),
	}

	go chat.writer()

	return chat
}

// Register registers new connection as a User.
func (c *Collective) Register(conn net.Conn) *Collaborator {
	user := &Collaborator{
		chat: c,
		conn: conn,
	}

	c.mu.Lock()
	{
		user.id = c.seq
		user.name = fmt.Sprintf("%d", c.seq) // TODO name from user...

		c.us = append(c.us, user)
		c.ns[user.name] = user

		c.seq++
	}
	c.mu.Unlock()

	fmt.Println("new user")

	err := user.writeNotice("hello", Object{
		"name": user.name,
	})
	if err != nil {
		log.Printf("error write notice: %v\n", err)
	}
	err = c.Broadcast("greet", Object{
		"name": user.name,
		"time": timestamp(),
	})
	if err != nil {
		log.Printf("error write notice: %v\n", err)
	}
	return user
}

// Remove removes user from chat.
func (c *Collective) Remove(user *Collaborator) {
	c.mu.Lock()
	removed := c.remove(user)
	c.mu.Unlock()

	if !removed {
		return
	}

	c.Broadcast("goodbye", Object{
		"name": user.name,
		"time": timestamp(),
	})
}

// Rename renames user.
func (c *Collective) Rename(user *Collaborator, name string) (prev string, ok bool) {
	c.mu.Lock()
	{
		if _, has := c.ns[name]; !has {
			ok = true
			prev, user.name = user.name, name
			delete(c.ns, prev)
			c.ns[name] = user
		}
	}
	c.mu.Unlock()

	return prev, ok
}

// Broadcast sends message to all alive users.
func (c *Collective) Broadcast(method string, params Object) error {
	var buf bytes.Buffer

	w := wsutil.NewWriter(&buf, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	r := Request{Method: method, Params: params}
	if err := encoder.Encode(r); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	c.out <- buf.Bytes()

	return nil
}

// writer writes broadcast messages from chat.out channel.
func (c *Collective) writer() {
	for bts := range c.out {
		c.mu.RLock()
		us := c.us
		c.mu.RUnlock()

		for _, u := range us {
			u := u // For closure.
			c.pool.Schedule(func() {
				u.writeRaw(bts)
			})
		}
	}
}

// mutex must be held.
func (c *Collective) remove(user *Collaborator) bool {
	if _, has := c.ns[user.name]; !has {
		return false
	}

	delete(c.ns, user.name)

	i := sort.Search(len(c.us), func(i int) bool {
		return c.us[i].id >= user.id
	})
	if i >= len(c.us) {
		panic("chat: inconsistent state")
	}

	without := make([]*Collaborator, len(c.us)-1)
	copy(without[:i], c.us[:i])
	copy(without[i:], c.us[i+1:])
	c.us = without

	return true
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// deadliner is a wrapper around net.Conn that sets read/write deadlines before
// every Read() or Write() call.
type deadliner struct {
	net.Conn
	t time.Duration
}

func (d deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}
