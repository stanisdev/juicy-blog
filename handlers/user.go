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
  if c.User.Authorized() {
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
  if c.User.Authorized() {
    http.Redirect(w, r, "/", 302)
    return
  }
  r.ParseForm()
  data := r.Form
  _, email := data["email"]
  _, password := data["password"]
  if !email || len(r.FormValue("email")) < 1 || !password || len(r.FormValue("password")) < 1 {
    c.SetFlash("Email or Password was not specified", "danger")
    http.Redirect(w, r, "/login", 302)
    return
  }
  var user models.User
  c.DB.First(&user, "email = ?", data["email"])
  if user.ID < 1 || !user.ComparePassword(r.FormValue("password")) {
    c.SetFlash("Incorrect Email or Password", "danger")
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
 * Save user's settings
 */
func UserSettingsSave(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  r.ParseForm()
  var user models.User

  if hasError, message := services.ValidateModel(&user, r.PostForm); hasError == true {
    c.SetFlash(message, "danger")
  } else {
    fmt.Println(c.User.Name)
    user.ID = c.User.ID
    c.DB.Save(&user)
    c.SetFlash("Settings was updated", "info")
  }
  http.Redirect(w, r, "/user/settings", 302)
}

/**
 * Change user's password
 */
func UserPasswordChange(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  r.ParseForm()
  password := r.FormValue("password")
  confirm := r.FormValue("password_confirm")
  if len(password) < 1 {
    c.SetFlash("Password value is empty", "danger")
    http.Redirect(w, r, "/user/settings", 302)
    return
  }
  if password != confirm {
    c.SetFlash("Password and Confirm do not match", "danger")
    http.Redirect(w, r, "/user/settings", 302)
    return
  }
  c.User.ChangePassword(password, c.DB)
  c.SetFlash("Password was changed", "info")
  http.Redirect(w, r, "/user/settings", 302)
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
