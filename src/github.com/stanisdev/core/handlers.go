package core

import (
  "fmt"
  "net/http"
  "html/template"
  "time"
  "github.com/stanisdev/helpers"
)

const viewPath = "src/github.com/stanisdev/templates/";

type Page struct {
    Title string
}

func loadTemplate(templateName string, w http.ResponseWriter, p *Page)  {
  w.Header().Set("Content-type", "text/html")
  t, err := template.ParseFiles(viewPath + templateName + ".html")
  if err != nil {
    fmt.Fprintf(w, "Template cannot be loaded")
  }
  t.Execute(w, p)
}

func Index(w http.ResponseWriter, r *http.Request, c *helpers.Containers) {

  defer c.DB.Close()
  // sessId := "cookie-sess-id" // генерировать случайную строку
  // db.CreateSession(sessId, "")
  //
  // expiration := time.Now().Add(365 * 24 * time.Hour)
  // cookie := http.Cookie{Name: "gcid", Value: sessId, Expires: expiration}
  // http.SetCookie(w, &cookie)

  p := &Page{Title: "My Blog"}
  loadTemplate("index", w, p)
}


type Cookie struct {
    Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string
    MaxAge   int
    Secure   bool
    HttpOnly bool
    Raw      string
    Unparsed []string // Raw text of unparsed attribute-value pairs
}

func AddArticle(w http.ResponseWriter, r *http.Request, c *helpers.Containers)  {
  // cookie, _ := r.Cookie("username")
  // for _, cookie := range r.Cookies() {
  //   fmt.Fprint(w, cookie.Name)
  // }
  //p := &Page{Title: "New Article"}
  ///loadTemplate("new-article", w, p)
}
