package core

import (
  "net/http"
  "strconv"
  m "github.com/stanisdev/models"
)

func Index(w http.ResponseWriter, r *http.Request, c *Containers) {
  c.Page.Title = "Mutual Blog"
}

func Login(w http.ResponseWriter, r *http.Request, c *Containers) {
  c.Page.Title = "Login to Blog"
  c.Page.Data["name"] = "John"
}

func Logout(w http.ResponseWriter, r *http.Request, c *Containers) {
  c.Session.Unset("user")
  http.Redirect(w, r, "/login", 302)
}

func LoginPost(w http.ResponseWriter, r *http.Request, c *Containers) {
  r.ParseForm()
  data := r.Form
  _, email := data["email"];
  _, password := data["password"];
  if !email || len(r.FormValue("email")) < 1 || !password || len(r.FormValue("password")) < 1 {
    c.SetFlash("Email or Password was not specified")
    http.Redirect(w, r, "/login", 302)
    return
  }
  var user m.User
  c.DB.First(&user, "email = ?", data["email"])
  if user.ID < 1 || !user.ComparePassword(r.FormValue("password")) {
    c.SetFlash("Incorrect Email or Password")
    http.Redirect(w, r, "/login", 302)
    return
  }
  c.Session.Set("user", strconv.Itoa(int(user.ID)))
  http.Redirect(w, r, "/", 302)
}

func Articles(w http.ResponseWriter, r *http.Request, c *Containers)  {
  c.Page.Title = "Articles"
}

func NewArticle(w http.ResponseWriter, r *http.Request, c *Containers)  {
  c.Page.Title = "New Article"
}

func NewArticlePost(w http.ResponseWriter, r *http.Request, c *Containers)  {
  
}