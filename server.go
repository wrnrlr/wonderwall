package main; import("crypto/rand"
"encoding/json";"errors";"fmt";"golang.org/x/crypto/bcrypt";"html/template";"log";"net/http";"regexp";"time")
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var (emailErr = errors.New("invalid email"); passwordErr = errors.New("invalid password"))
type Email string; func (e Email) valid() bool { return len(e) > 3 && len(e) < 255 && emailRegex.MatchString(string(e)) }
type Password string; func (p Password) valid() bool { return len(p) > 8 && len(p) < 255 }
func (p Password) HashPassword() (PasswordHash, error) { return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost) }
type RegistrationForm struct { Email Email `json:"email"`; Password Password `json:"password"`}
func (f RegistrationForm) validate() error { if !f.Email.valid() { return emailErr } else if len(f.Password) < 8 { return passwordErr }; return nil}
type EmailForm struct { Email Email }; func (f EmailForm) validate() error { if !f.Email.valid() { return emailErr }; return nil}
type PasswordForm struct { Password Password }; func (f PasswordForm) validate() error { if !f.Password.valid() { return passwordErr }; return nil}
type PasswordHash []byte; func (h PasswordHash) ComparePassword(p Password) error { return bcrypt.CompareHashAndPassword(h, []byte(p)) }
type Registration struct { ID Token; Email Email; PasswordHash PasswordHash; CreatedAt time.Time; VerifiedAt *time.Time }
func (r *Registration) Key() Key { if r == nil { return Key("registration:")} else { return append([]byte("registration:"),r.ID...) } }
func (r *Registration) Serialize() ([]byte, error) { return serialize(r) }
func (r *Registration) Deserialize(b []byte) error { return deserialize(b, r) }
type User struct { ID string; Email Email; PasswordHash PasswordHash }
func (u *User) Key() Key { if u == nil { return Key("user:")} else { return Key("user:"+u.ID) } }
type Wall struct { ID string; Elements []interface{}}
func (w *Wall) Key() Key { if w == nil { return Key("wall:")} else { return Key("wall:"+w.ID) } }
type FindUserByEmail interface { FindUserByEmail(*Txn, Email) (*User, error) }
type FindUserById interface { FindUserById(*Txn, string) (*User, error) }
type CreateUser interface { CreateUser(*Txn, *User) error }
type UpdateUser interface { UpdateUser(*Txn, *User) error }
type DeleteUser interface { DeleteUser(*Txn, *User) error }
type UserService interface { CreateUser; UpdateUser; DeleteUser; FindUserById; FindUserByEmail }
type Users struct { users []*Users }
func (s Users) CreateUser(*Txn, *User) error { return nil }
func (s Users) UpdateUser(*Txn, *User) error { return nil }
func (s Users) DeleteUser(*Txn, *User) error { return nil }
func (s Users) FindUserById(*Txn, string) (*User, error) { return nil, nil }
func (s Users) FindUserByEmail(*Txn, Email) (*User, error) { return nil, nil }
var ( applicationJson = "application/json")
func ContentType(t string, w http.ResponseWriter, r *http.Request) bool { if r.Header.Get("Content-Type") != t { w.WriteHeader(http.StatusUnsupportedMediaType); return false } else { return true } }
func writeTmpl(w http.ResponseWriter, name string, i interface{}) {
	indexTmpl, err := template.ParseFiles(fmt.Sprintf("./template/%s.html", name)); if err != nil {panic(err)}
	if err = indexTmpl.Execute(w, nil); err != nil {panic(err)}}
func writeError(w http.ResponseWriter, err error) { writeTmpl(w, "500", err); w.WriteHeader(500) }
func RenderTemplate(name string) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		writeTmpl(w, name, nil)
	}
}
type Token []byte
func GenerateToken(n int) (Token,error) {
	b := make([]byte, n); _, err := rand.Read(b)
	if err != nil {return nil, err}; return b, nil}
type CreateRegistration interface { CreateRegistration(*Txn, Email, Password) (*Registration,error) }
type RegistrationService interface { CreateUser; FindUserByEmail; }
type Registrations struct { DB *Store }
func (s Registrations) CreateRegistration(txn *Txn, email Email, password Password) (*Registration,error) {
	passwordHash, err := password.HashPassword(); if err != nil {return nil, err}
	id, err := GenerateToken(32); if err != nil {return nil, err}
	now := time.Now()
	r := &Registration{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt:now}
	err = s.DB.Set(txn, r); if err != nil {return nil, err}
	return r, nil
}
type Emails struct {}
func (s Emails) SendEmail(name string, i interface{}) error { return nil }
var duplicateEmailErr = errors.New("duplicate email")
func PostRegistration(db *Store, registrations CreateRegistration, users FindUserByEmail, emails *Emails) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var (f RegistrationForm; u *User; reg *Registration; err error)
		if !ContentType(applicationJson, w, r) { return }
    if err := json.NewDecoder(r.Body).Decode(&f); err != nil { w.WriteHeader(http.StatusBadRequest);return }
    if err := f.validate(); err != nil { w.WriteHeader(http.StatusBadRequest);return }
    err = db.Update(func(txn*Txn)error{
    	u, err = users.FindUserByEmail(txn, f.Email); if err != nil { return err } else if u != nil { return duplicateEmailErr }
    	reg, err = registrations.CreateRegistration(txn, f.Email, f.Password); if err != nil { return err }
    	return nil
    })
    if err != nil { w.WriteHeader(http.StatusBadRequest);return }
    emails.SendEmail("verify-email", reg.ID)
	}
}
func RegistrationHandler(get, post http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { get(w, r) } else
		if r.Method == "POST" { post(w, r) }
	}
}
func main() {
	indexHandler := func(w http.ResponseWriter, r *http.Request) { writeTmpl(w, "index", nil) }
	sandboxHandler := func(w http.ResponseWriter, r *http.Request) { writeTmpl(w, "sandbox", nil) }
	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { writeTmpl(w, "login", nil); return }
		if r.Method == "POST" {
			var form Registration; if err := json.NewDecoder(r.Body).Decode(&form); err != nil { writeError(w, err); return }
			http.SetCookie(w, &http.Cookie{Name:"session",Value:"session"})}}
	logoutHandler := func(w http.ResponseWriter, r *http.Request) {}
	forgotPasswordHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { writeTmpl(w, "forgot-password", nil); return }
		if r.Method == "POST" {
			var form EmailForm; if err := json.NewDecoder(r.Body).Decode(&form); err != nil { writeError(w, err); return }
			/*TODO*/}}
	termsHandler := func(w http.ResponseWriter, r *http.Request) { writeTmpl(w, "terms", nil) }
	store := &Store{}
	users := &Users{}
	registrations := &Registrations{}
	emails := &Emails{}
	postRegistration := PostRegistration(store, registrations, users, emails)
	getRegistration := RenderTemplate("join")
	wrapper := noCacheMiddleware
	http.HandleFunc("/", wrapper(indexHandler))
	http.HandleFunc("/sandbox", wrapper(sandboxHandler))
	http.HandleFunc("/join", wrapper(RegistrationHandler(getRegistration, postRegistration)))
	http.HandleFunc("/terms", wrapper(termsHandler))
	http.HandleFunc("/login", wrapper(loginHandler))
	http.HandleFunc("/logout", wrapper(logoutHandler))
	http.HandleFunc("/forgot-password", wrapper(forgotPasswordHandler))
	http.HandleFunc("/reset-password", wrapper(loginHandler))
	http.HandleFunc("/walls", wrapper(loginHandler))
	http.HandleFunc("/wall", wrapper(loginHandler))
	http.Handle("/static/", wrapper(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP))
	log.Fatal(http.ListenAndServe(":9999", nil))
}
