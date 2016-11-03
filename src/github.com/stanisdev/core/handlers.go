package core

import (
  "net/http"
  "fmt"
  m "github.com/stanisdev/models"
)

func Index(w http.ResponseWriter, r *http.Request, c *Containers) {
  defer c.DB.Close()
  // c.Session.Set("city", "Tokio")
  city, _ := c.Session.Get("city")
  p := &Page{Title: "My Blog " + city}
  loadTemplate("index", w, p)
}

func Login(w http.ResponseWriter, r *http.Request, c *Containers) {
  c.Page.Title = "Login Super 22"
  c.Page.Data["name"] = "John"
  fmt.Println(c.GetFlash())
}

func LoginPost(w http.ResponseWriter, r *http.Request, c *Containers) {
  r.ParseForm()
  data := r.Form
  if _, email := data["email"]; !email {
    c.SetFlash("my_message 2")
    fmt.Fprint(w, "Incorrect login/password")
    return
  }
  if _, password := data["password"]; !password {
    c.SetFlash("my_message 2")
    fmt.Fprint(w, "Incorrect login/password")
    return
  }
  var user m.User
  c.DB.First(&user, "email = ?", data["email"])
  if user.ID < 1 || !user.ComparePassword(r.FormValue("password")) {
    c.SetFlash("my_message 2")
    http.Redirect(w, r, "/login", 302)
    return
  }
}

func AddArticle(w http.ResponseWriter, r *http.Request, c *Containers)  {
  //p := &Page{Title: "New Article"}
  ///loadTemplate("new-article", w, p)
}
