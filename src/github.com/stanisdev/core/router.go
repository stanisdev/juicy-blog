package core

import (
  "net/http"
  "fmt"
  "github.com/stanisdev/db"
  "runtime"
  "reflect"
  "strings"
)

type RouterHandler func(http.ResponseWriter, *http.Request, *Containers)

type Router struct {
  Handlers map[string]map[string]RouterHandler
}

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
  hs, exists := self.Handlers[r.URL.Path]
  if !exists {
    self.notFound(w)
    return
  }
  h, exists := hs[r.Method]
  if !exists {
    self.notFound(w)
    return
  }
  // Handler has found in map
  dbConnection := db.Connect()
  c := &Containers{DB: dbConnection, Session: &SessionManager{DB: dbConnection}, Page: &Page{}}
  c.Page.Data = make(map[string]interface{})
  c.Session.Start(w, r)
  methodName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
  methodName = methodName[strings.LastIndex(methodName, ".")+1:]
  h(w, r, c)
  loadTemplate(strings.ToLower(methodName), w, c.Page)
}

func (self *Router) Start() {
  http.HandleFunc("/", self.handler)
  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}