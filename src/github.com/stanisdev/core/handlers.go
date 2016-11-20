package core

import (
  "net/http"
  "strconv"
  "github.com/stanisdev/models"
  "fmt"
  "math"
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
  var user models.User
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
      c.BadRequest = "Page parameter must be a number"
      return
    }
    if p > 1 {
      offset = (p - 1) * 5
    }
  }
  var count int
  c.DB.Model(&models.Article{}).Count(&count)
  c.Page.Data["count"] = count
  c.Page.Title = "Articles"
  if count < 1 {
    return
  }
  if count <= offset {
    c.BadRequest = "Page does not exist"
    return
  }
  if count > 5 {
    pagesCount := math.Ceil(float64(count) / 5)
    pg, _ := strconv.Atoi(page)
    if pg < 1 {
      pg = 1
    }
    c.Page.Data["pagination"] = makePagination(pg, int(pagesCount))
  }
  c.Page.Data["articles"] = c.Models.GetArticles(5, offset)
}

/**
 * Create new Article
 */
func ArticleNew(w http.ResponseWriter, r *http.Request, c *Containers)  {
  c.Page.Title = "New Article"
}

/**
 * Create new Article (POST)
 */
func ArticleNewPost(w http.ResponseWriter, r *http.Request, c *Containers)  {
  r.ParseForm()
  var article models.Article
  if hasError, message := ValidateModel(&article, r.PostForm); hasError == true {
    c.SetFlash(message)
  } else {
    article.UserID = c.Page.User.ID
    c.DB.Create(&article)
    if article.ID < 1 {
      c.SetFlash("Article cannot be created")
    }
  }
  fmt.Println("Posted")
  http.Redirect(w, r, "/articles/new", 302)
}

/**
 * View Article (GET)
 */
func ArticleView(w http.ResponseWriter, r *http.Request, c *Containers)  {
  // @TODO check to number (create separate function, that check whether number)
  id, _ := strconv.Atoi(c.Params["id"])
  var article models.Article
  c.DB.Find(&article, id)
  title := "Article"
  if article.ID < 1 {
    c.Page.Data["notFound"] = true
  } else {
    title += " :: " + article.Title
    c.Page.Data["article"] = article
  }
  c.Page.Title = title
}