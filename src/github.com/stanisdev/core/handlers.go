package core

import (
  "fmt"
  "net/http"
  "html/template"
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

func Index(w http.ResponseWriter, r *http.Request, c *Containers) {
  defer c.DB.Close()
  p := &Page{Title: "My Blog"}
  loadTemplate("index", w, p)
}

func AddArticle(w http.ResponseWriter, r *http.Request, c *Containers)  {
  //p := &Page{Title: "New Article"}
  ///loadTemplate("new-article", w, p)
}
