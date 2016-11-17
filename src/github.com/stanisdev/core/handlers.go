package core

import (
  "net/http"
  "strconv"
  m "github.com/stanisdev/models"
)

/**
 * Index page
 */
func Index(w http.ResponseWriter, r *http.Request, c *Containers) {
  c.Page.Title = "Mutual Blog"
}

/**
 * Login page
 */
func Login(w http.ResponseWriter, r *http.Request, c *Containers) {
  if c.Page.User.Authorized() {
    http.Redirect(w, r, "/", 302)
    return
  }
  c.Page.Title = "Login to Blog"
  c.Page.Data["name"] = "John"
}

/**
 * Login (POST)
 */
func LoginPost(w http.ResponseWriter, r *http.Request, c *Containers) {
  if c.Page.User.Authorized() {
    http.Redirect(w, r, "/", 302)
    return
  }
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

/**
 * Logout page
 */
func Logout(w http.ResponseWriter, r *http.Request, c *Containers) {
  c.Session.Unset("user")
  http.Redirect(w, r, "/login", 302)
}

/**
 * Articles list
 */
func Articles(w http.ResponseWriter, r *http.Request, c *Containers)  {
  var offset int = 0
  var page string = c.Params["page"]
  if len(page) > 0 { // Page param exists
    p, err := strconv.Atoi(page)
    if err != nil {
      http.Redirect(w, r, "/articles", 302)
      return
    }
    offset = p
  }
  c.Page.Data["articles"] = c.Models.GetArticles(0, offset)
  c.Page.Title = "Articles"
}

/**
 * Create new Article
 */
func NewArticle(w http.ResponseWriter, r *http.Request, c *Containers)  {
  c.Page.Title = "New Article"
}

/**
 * Create new Article (POST)
 */
func NewArticlePost(w http.ResponseWriter, r *http.Request, c *Containers)  {
  r.ParseForm()
  var article m.Article
  if hasError, message := ValidateModel(&article, r.PostForm); hasError == true {
    c.SetFlash(message)
  } else {
    article.UserID = c.Page.User.ID
    c.DB.Create(&article)
    if article.ID < 1 {
      c.SetFlash("Article cannot be created")
    }
  }
  http.Redirect(w, r, "/articles/new", 302)
}