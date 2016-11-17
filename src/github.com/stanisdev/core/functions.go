package core

import (
  "net/http"
  "time"
  "math/rand"
  "html/template"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/url"
  "github.com/gorilla/schema"
  validator "github.com/asaskevich/govalidator"
  "strings"
)

const viewPath = "src/github.com/stanisdev/templates/"

/**
 * Container's methods
 */
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
    c.DB.Select("id, name, email").First(&c.Page.User, userId)
  }
}

/**
 * General Functions
 */
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
  tplFuncMap := make(template.FuncMap)
  tplFuncMap["IsArticles"] = func (url string) bool {
    return len(url) > 8 && url[:9] == "/articles"
  }
  t, err := template.New("").Funcs(tplFuncMap).ParseFiles(viewPath + "/layouts/layout.html", viewPath + templateName + ".html")

  if err != nil {
    fmt.Fprintf(w, "Template cannot be loaded")
  }
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

func MakePagination(currentPage int, pageCount int) {
  var marks []struct{Title string; Clickable bool; Number interface{}}
  fmt.Println(marks)
}