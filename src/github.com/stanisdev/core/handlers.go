package core

import (
  "net/http"
  "fmt"
)

func Index(w http.ResponseWriter, r *http.Request, c *Containers) {
  defer c.DB.Close()
  // c.Session.Set("city", "Tokio")
  city, _ := c.Session.Get("city")
  p := &Page{Title: "My Blog " + city}
  loadTemplate("index", w, p)
}

func Login(w http.ResponseWriter, r *http.Request, c *Containers) {
  if r.Method == "POST" {
    r.ParseForm()
    fmt.Println(r.Form)
    fmt.Println(r.FormValue("email"))
  }
  p := &Page{Title: "Login"}
  loadTemplate("login", w, p)
}

func AddArticle(w http.ResponseWriter, r *http.Request, c *Containers)  {
  //p := &Page{Title: "New Article"}
  ///loadTemplate("new-article", w, p)
}
