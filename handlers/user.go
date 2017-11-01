package handlers

import (
  "net/http"
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/models"
  "strconv"
  "fmt"
)

/**
 * Login page
 */
func Login(w http.ResponseWriter, r *http.Request, c *services.Containers) {
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
func LoginPost(w http.ResponseWriter, r *http.Request, c *services.Containers) {
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
func Logout(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  c.Session.Unset("user")
  fmt.Println("Success")
  http.Redirect(w, r, "/login", 302)
}

/**
 * User settings
 */
func UserSettings(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  c.Page.Title = "User settings"
}

/**
 * View Profile
 */
func UserView(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var user models.User
  c.DB.Find(&user, id)

  if user.ID < 1 {
    c.Page.Data["notFound"] = true
  } else {
    c.Page.Data["user"] = user
  }
  c.Page.Title = "View user profile"
}
