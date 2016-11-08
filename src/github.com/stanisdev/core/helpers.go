package core

import (
  "github.com/jinzhu/gorm"
  "net/http"
  "time"
  "math/rand"
  "html/template"
  "fmt"
  m "github.com/stanisdev/models"
  "encoding/json"
  "io/ioutil"
  "net/url"
  "github.com/gorilla/schema"
  validator "github.com/asaskevich/govalidator"
  "strings"
)

func ValidateModel(modelInstance interface{}, formData url.Values) (bool, string) {
  decoder := schema.NewDecoder()
  if err := decoder.Decode(modelInstance, formData); err != nil {
    return true, "Fields cannot be parsed"
  }
  if _, err := validator.ValidateStruct(modelInstance); err != nil {
    var message string = err.Error()
    var splited []string = strings.Split(message[:len(message)-1], ";")
    return true, strings.Join(splited, "<br/>")
  } else {
    return false, ""
  }
}

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

type Config struct {
  DbName string `json:"db_name"`
  DbUser string `json:"db_user"`
  DbPass string `json:"db_pass"`
  UrlsWithoutTemplate []string `json:"urls_without_template"`
  ProtectedUrls []struct{
    Url string 
    Method string
  } `json:"protected_urls"`
}

func GetConfig() *Config {
  raw, err := ioutil.ReadFile("./config.json")
  if err != nil {
    panic("Config file cannot be loaded")
  }
  var config Config
  if err := json.Unmarshal(raw, &config); err != nil {
    panic("JSON config cannot be parsed")
  }
  return &config
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