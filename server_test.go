package wonderwall

import (
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPassword(t *testing.T) {
	t.Run("Valid password has minimum length 8", func(t *testing.T) { assert.False(t, Password("short").valid()) })
}

func TestEmail(t *testing.T) {
	assert.False(t, Email("short").valid())
	assert.False(t, Email("short@").valid())
}

func TestAuthForm(t *testing.T) {
	assert.Equal(t, RegistrationForm{"test@example.com", "short"}.validate(), passwordErr)
	assert.Equal(t, RegistrationForm{"test", "Password12345!"}.validate(), emailErr)
}

var (
	aliceEmail           = Email("alice@example.com")
	alicePassword        = Password([]byte("Password1234!"))
	alicePasswordHash, _ = alicePassword.HashPassword()
)
var alice = User{xid.New(), aliceEmail, alicePasswordHash, "alice"}

type findUserByEmailReturnUser struct{}

func (s findUserByEmailReturnUser) FindUserByEmail(email Email) (*User, error) {
	if email == "alice@example.com" {
		return &alice, nil
	}
	return nil, nil
}

type createUser struct{}

func (c createUser) CreateUser(u *User) error { return nil }

type registrationServiceMock struct {
	findUserByEmailReturnUser
	createUser
}

//func TestRegistrationHandler(t *testing.T) {
//	t.Run("GET /join", func(t *testing.T) {
//		h := GetPostRouter(struct{findUserByEmailReturnUser;createUser}{})
//		req, err := http.NewRequest("GET", "/join", nil)
//		if err != nil { t.Fatal(err) }
//		w := httptest.NewRecorder()
//		h.ServeHTTP(w, req)
//		resp := w.Result()
//		assert.Equal(t, resp.StatusCode, 200)
//	})
//	t.Run("POST valid form", func(t *testing.T) {
//		h := GetPostRouter(struct{findUserByEmailReturnUser;createUser}{})
//		f := Registration{aliceEmail, alicePassword}
//		b, _ := json.Marshal(&f)
//		r, err := http.NewRequest("POST", "/join", bytes.NewReader(b))
//		if err != nil { t.Fatal(err) }
//		r.Header.Set("Content-Type", "application/json")
//		w := httptest.NewRecorder()
//		h.ServeHTTP(w, r)
//		resp := w.Result()
//		assert.Equal(t, resp.StatusCode, 200)
//	})
//}

//func TestLoginHandler(t *testing.T) {
//	t.Run("GET /login", func(t *testing.T) {
//		h := PostLogin(struct{findUserByEmailReturnUser;createUser}{})
//		req, err := http.NewRequest("GET", "/login", nil)
//		if err != nil { t.Fatal(err) }
//		w := httptest.NewRecorder()
//		h.ServeHTTP(w, req)
//		resp := w.Result()
//		assert.Equal(t, resp.StatusCode, 200)
//	})
//	t.Run("POST valid form", func(t *testing.T) {
//		h := GetPostRouter(struct{findUserByEmailReturnUser;createUser}{})
//		f := Registration{"alice@example.com", "password1234!"}
//		b, _ := json.Marshal(&f)
//		r, err := http.NewRequest("POST", "/login", bytes.NewReader(b))
//		if err != nil { t.Fatal(err) }
//		r.Header.Set("Content-Type", "application/json")
//		w := httptest.NewRecorder()
//		h.ServeHTTP(w, r)
//		resp := w.Result()
//		assert.Equal(t, resp.StatusCode, 200)
//	})
//	t.Run("POST unknown email", func(t *testing.T) {
//		h := GetPostRouter(struct{findUserByEmailReturnUser;createUser}{})
//		f := Registration{"bob@example.com", "Password1234!"}
//		b, _ := json.Marshal(&f)
//		r, err := http.NewRequest("POST", "/login", bytes.NewReader(b))
//		if err != nil { t.Fatal(err) }
//		r.Header.Set("Content-Type", "application/json")
//		w := httptest.NewRecorder()
//		h.ServeHTTP(w, r)
//		resp := w.Result()
//		assert.Equal(t, resp.StatusCode, 400)
//	})
//	t.Run("POST invalid password", func(t *testing.T) {
//		h := GetPostRouter(struct{findUserByEmailReturnUser;createUser}{})
//		f := Registration{aliceEmail, "NotThePassword1234!"}
//		b, _ := json.Marshal(&f)
//		r, err := http.NewRequest("POST", "/login", bytes.NewReader(b))
//		if err != nil { t.Fatal(err) }
//		r.Header.Set("Content-Type", "application/json")
//		w := httptest.NewRecorder()
//		h.ServeHTTP(w, r)
//		resp := w.Result()
//		assert.Equal(t, resp.StatusCode, 400)
//	})
//}
