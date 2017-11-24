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
    c.Redirect("/")
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
    c.Redirect("/")
    return
  }
  r.ParseForm()
  data := r.Form
  _, email := data["email"]
  _, password := data["password"]
  if !email || len(r.FormValue("email")) < 1 || !password || len(r.FormValue("password")) < 1 {
    c.SetFlash("Email or Password was not specified", "danger")
    c.Redirect("/login")
    return
  }
  var user models.User
  c.DB.First(&user, "email = ?", data["email"])
  if user.ID < 1 || !user.ComparePassword(r.FormValue("password")) {
    c.SetFlash("Incorrect Email or Password", "danger")
    c.Redirect("/login")
    return
  }
  c.Session.Set("user", strconv.Itoa(int(user.ID)))
  c.Redirect("/")
}

/**
 * Logout page
 */
func Logout(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  c.NoTemplate = true
  c.Session.Unset("user")
  c.Redirect("/login")
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
  c.Redirect("/user/settings")
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
  } else if password != confirm {
    c.SetFlash("Password and Confirm do not match", "danger")
  } else {
    c.User.ChangePassword(password, c.DB)
    c.SetFlash("Password was changed", "info")
  }
  c.Redirect("/user/settings")
}

/**
 * Subscribe to User's articles
 */
func UserSubscribing(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  userId, _ := c.IntUrlParams["id"]
  if uint(userId) == c.User.ID {
    c.SetFlash("It's disallow to subscribe to myself", "danger")
    c.Redirect("/")
    return
  }
  var user models.User
  c.DB.Find(&user, userId) // Вынести в middleware
  c.Page.Title = "Subscribe to user"

  if user.ID < 1 {
    return
  }
  subscribingPage := "/user/" + strconv.Itoa(int(userId))
  if c.Models.SubscriberExists(uint(userId), c.User.ID) {
    c.DB.Where("user_id = ? AND subscriber_id = ?", userId, c.User.ID).Delete(models.Subscriber{})
    c.SetFlash("You unsubscribed successfully", "info")
    c.Redirect(subscribingPage)
    return
  }
  newRecord := models.Subscriber{ UserID: uint(userId), SubscriberID: c.User.ID }
  c.DB.Create(&newRecord)
  c.SetFlash("You subscribed successfully", "info")
  c.Redirect(subscribingPage)
}

/**
 * View Profile
 */
func UserView(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id, _ := c.IntUrlParams["id"]
  var user models.User
  c.DB.Find(&user, id)

  if user.ID < 1 {
    c.Page.Data["notFound"] = true
  } else {
    c.Page.Data["user"] = user
  }
  c.Page.Title = "View user profile"
  c.Page.Data["isSubscribed"] = c.Models.SubscriberExists(uint(id), c.User.ID)
}
