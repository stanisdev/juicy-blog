package core

import (
  "github.com/jinzhu/gorm"
  "net/http"
  "time"
  "math/rand"
  "github.com/stanisdev/db"
)

type Containers struct {
  DB *gorm.DB
  Session *SessionManager
}

type Cookie struct {
    Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string
    MaxAge     int
    Secure     bool
    HttpOnly   bool
    Raw        string
    Unparsed   []string // Raw text of unparsed attribute-value pairs
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *Containers)) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request)  {
    dbConnection := db.Connect()
    c := &Containers{DB: dbConnection, Session: &SessionManager{DB: dbConnection}}
    c.Session.Start(w, r)
    fn(w, r, c)
  }
}

func GenerateRandomString(l int) string {
  rand.Seed(time.Now().UTC().UnixNano())
  bytes := make([]byte, l)
  for i := 0; i < l; i++ {
      bytes[i] = byte(randInt(65, 90))
  }
  return string(bytes)
}

func randInt(min int, max int) int {
  return min + rand.Intn(max-min)
}