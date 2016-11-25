package core

import (
  "net/http"
  "fmt"
  "github.com/stanisdev/db"
  "github.com/stanisdev/models"
  "runtime"
  "reflect"
  "strings"
  "regexp"
  "html/template"
)

type RouterHandler func(http.ResponseWriter, *http.Request, *Containers)

func (self *Router) defineHandler(url string) {
  if len(self.Handlers[url]) < 1 {
    self.Handlers[url] = make(map[string]RouterHandler)
  }
}

func (self *Router) GET(url string, fn RouterHandler) {
  self.defineHandler(url)
  self.Handlers[url]["GET"] = fn
}

func (self *Router) POST(url string, fn RouterHandler) {
  self.defineHandler(url)
  self.Handlers[url]["POST"] = fn
}

func (self *Router) notFound(w http.ResponseWriter) {
  w.WriteHeader(http.StatusNotFound)
  fmt.Fprint(w, "Page not found 404")
}

func (self *Router) handler(w http.ResponseWriter, r *http.Request) {
  url := r.URL.Path 
  if url == "/favicon.ico" {
    return
  }
  reqData := make(map[string]string)
  hs, exists := self.Handlers[url]
  if !exists { // First compare by simple equal
    var matches []string
    reg, _ := regexp.Compile("\\/:([a-z]+\\/??)") 
    for pattern, _ := range self.Handlers { // Compare to patterns
      if strings.Contains(pattern, ":") {
        changedPattern := reg.ReplaceAllString(pattern, "(\\/[^\\/]+)")
        re, _ := regexp.Compile("^" + changedPattern + "$") 
        matches = re.FindStringSubmatch(url)
        if len(matches) > 0 { // Pattern has found
          hs, _ = self.Handlers[pattern]
          params := reg.FindAllString(pattern, len(matches)-1)
          for key, val := range params {
            var reqValue string = matches[key+1]
            if len(reqValue) > 0 {
              reqValue = reqValue[1:]
            }
            reqData[val[2:]] = reqValue
          }
          break
        }
      }
    }
    if len(matches) < 1 {
      self.notFound(w)
      return
    }
  }
  h, exists := hs[r.Method]
  if !exists {
    self.notFound(w)
    return
  }
  // Handler has found in map
  dbConnection := db.Connect(self.Config.DbUser, self.Config.DbPass, self.Config.DbName)

  c := &Containers{
    DB: dbConnection, 
    Models: models.StaticMethods{DB: dbConnection}, 
    Session: SessionManager{}, 
    Page: Page{},
    Params: reqData,
  }
  c.Page.Data = make(map[string]interface{})
  c.Session.Start(w, r)
  methodName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
  methodName = methodName[strings.LastIndex(methodName, ".")+1:]
  c.Auth()
  // Check restricted urls
  for _, value := range self.Config.ProtectedUrls {
    if value.Url == url && value.Method == r.Method && !c.Page.User.Authorized() {
      self.notFound(w)
      return
    }
  }
  h(w, r, c)
  if len(c.BadRequest) > 0 { // Some intolerable errors
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, c.BadRequest)
    return
  }
  if r.Method == "GET" {
    for _, value := range self.Config.UrlsWithoutTemplate {
      if value == url {
        return
      }
    }
    message, hasFlash := c.GetFlash()
    if (hasFlash == true) {
      c.Page.Flash = template.HTML(message)
    }
    c.Page.Url = url
    loadTemplate(strings.ToLower(methodName), w, &c.Page)
  }
}

func (self *Router) Start() {
  http.HandleFunc("/", self.handler)
  http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}