package session

import (
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"net/http"
	"os"
)

// Name is session name
const Name = "fascia.io"

// Session has session store
type Session struct {
	CookieStore *sessions.CookieStore
}

var sharedInstance = New()

// SharedInstance make a singleton instance
func SharedInstance() *Session {
	return sharedInstance
}

// New returns a Session struct
func New() *Session {
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET")))
	return &Session{
		CookieStore: store,
	}
}

// Get a value from session
func (u *Session) Get(r *http.Request, key string) (interface{}, error) {
	s, err := u.CookieStore.Get(r, Name)
	if err != nil {
		return nil, errors.Wrap(err, "cookie error")
	}
	v := s.Values[key]
	return v, nil
}

// Set a value to session
func (u *Session) Set(r *http.Request, w http.ResponseWriter, key string, value interface{}, option ...*sessions.Options) error {
	s, err := u.CookieStore.Get(r, Name)
	if err != nil {
		return errors.Wrap(err, "cookie error")
	}
	if len(option) > 0 {
		s.Options = option[0]
	}
	s.Values[key] = value
	return s.Save(r, w)
}

// Clear all session
func (u *Session) Clear(r *http.Request, w http.ResponseWriter) error {
	s, err := u.CookieStore.Get(r, Name)
	if err != nil {
		return errors.Wrap(err, "cookie error")
	}
	s.Options = &sessions.Options{
		MaxAge:   -1,
		HttpOnly: true,
	}
	return s.Save(r, w)
}
