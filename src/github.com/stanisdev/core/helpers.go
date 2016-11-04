package core

import (
  "github.com/jinzhu/gorm"
  "net/http"
  "time"
  "math/rand"
  "html/template"
  "fmt"
  m "github.com/stanisdev/models"
)

const viewPath = "src/github.com/stanisdev/templates/";

type Containers struct {
  DB *gorm.DB
  Session *SessionManager
  Page *Page
}

func (c *Containers) SetFlash(value string) {
  c.Session.Set("flash", value)
}

func (c *Containers) GetFlash() (string, bool) {
  value, exists := c.Session.Get("flash")
  if exists == true {
    c.Session.Unset("flash")
    return value, true
  } else {
    return "", false
  }
}

func (c *Containers) Auth() {
  userId, isAuth := c.Session.Get("user") 
  if isAuth == true {
    var user m.User
    c.DB.Select("id, name, email").First(&user, userId)
    if user.ID > 0 {
      c.Page.User = &user
    }
  }
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
  Flash string
  Url string
  User *m.User
  Data map[string]interface{}
}

func loadTemplate(templateName string, w http.ResponseWriter, p *Page)  {
  w.Header().Set("Content-type", "text/html")
  t, err := template.ParseFiles(viewPath + "/layouts/layout.html", viewPath + templateName + ".html")
  if err != nil {
    fmt.Fprintf(w, "Template cannot be loaded")
  }
  //t.Execute(w, p)
  t.ExecuteTemplate(w, "layout", p)
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