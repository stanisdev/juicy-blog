package core

import (
  "github.com/jinzhu/gorm"
  "net/http"
  "time"
  "math/rand"
  "html/template"
  "fmt"
  "github.com/stanisdev/db"
)

const viewPath = "src/github.com/stanisdev/templates/";

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

type Page struct {
    Title string
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *Containers)) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request)  {
    dbConnection := db.Connect()
    c := &Containers{DB: dbConnection, Session: &SessionManager{DB: dbConnection}}
    c.Session.Start(w, r)
    fn(w, r, c)
  }
}

func loadTemplate(templateName string, w http.ResponseWriter, p *Page)  {
  w.Header().Set("Content-type", "text/html")
  t, err := template.ParseFiles(viewPath + templateName + ".html")
  if err != nil {
    fmt.Fprintf(w, "Template cannot be loaded")
  }
  t.Execute(w, p)
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