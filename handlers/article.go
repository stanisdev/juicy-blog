package handlers

import (
  "net/http"
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/models"
  "sync"
  "math"
  "strconv"
  "time"
)

/**
 * Index page
 */
func Index(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  var articlesCount, usersCount int
  waitGroup := sync.WaitGroup{}
  waitGroup.Add(2)
  go func() {
    c.DB.Model(models.Article{}).Count(&articlesCount)
    waitGroup.Done()
  }()
  go func() {
    c.DB.Model(models.User{}).Count(&usersCount)
    waitGroup.Done()
  }()
  waitGroup.Wait()
  c.Page.Data["ArticlesCount"] = articlesCount
  c.Page.Data["UsersCount"] = usersCount
  c.Page.Title = "Mutual Blog"
}

/**
 * Articles list
 */
func Articles(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  page := c.GetParamByType(services.TypedRequestParam{ Name: "page", Type: "int", DefaultValue: 1 }).(int)
  if page < 1 {
    c.BadRequest = "Page must be greater then zero"
    return
  }
  var offset int = 0
  if page > 1 {
    offset = (page - 1) * 5
  }
  var count int
  c.DB.Model(models.Article{}).Count(&count)
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
    c.Page.Data["pagination"] = services.MakePagination(page, int(pagesCount))
  }
  articles := c.Models.GetArticles(5, offset)
  c.Page.Data["articles"] = articles
}

/**
 * Create new Article
 */
func ArticleNew(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  c.Page.Title = "New Article"
}

/**
 * Create new Article (POST)
 */
func ArticleNewPost(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  r.ParseForm()
  var article models.Article
  article.UserID = c.Page.User.ID
  if hasError, message := services.ValidateModel(&article, r.PostForm); hasError == true {
    c.SetFlash(message)
    http.Redirect(w, r, "/articles/new", 302)
  } else {
    if err := c.DB.Create(&article).Error; err != nil {
      c.SetFlash("Article cannot be created")
      http.Redirect(w, r, "/articles/new", 302)
    } else {
      http.Redirect(w, r, "/article/" + strconv.Itoa(int(article.ID)), 302)
    }
  }
}

/**
 * View Article (GET)
 */
func ArticleView(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var article struct{ID uint; Title string; Content string; CreatedAt time.Time; Userid uint; Username string}
  c.Models.FindArticleById(&article, id)

  if article.ID < 1 {
    c.Page.Title = "Article not found"
    c.Page.Data["notFound"] = true
  } else {
    c.Page.Title = article.Title
    c.Page.Data["article"] = article
    c.Page.Data["canEdit"] = article.Userid == c.Page.User.ID
  }
}

/**
 * Edit Article (GET)
 */
func ArticleEdit(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var article models.Article
  c.DB.Find(&article, id)

  if article.ID < 1 || article.UserID != c.Page.User.ID {
    c.Page.Title = "Not allowed to edit"
    c.Page.Data["isResctrictedEdit"] = true
  } else {
    c.Page.Title = "Editing :: " + article.Title
    c.Page.Data["article"] = article
  }
}

/**
 * Edit Article (POST)
 */
func ArticleEditPost(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var article models.Article
  c.DB.Find(&article, id)
  if article.ID < 1 || article.UserID != c.Page.User.ID {
    c.SetFlash("Not allowed to edit")
    http.Redirect(w, r, "/articles", 302)
    return
  }
  r.ParseForm()
  article.Title = r.FormValue("title")
  article.Content = r.FormValue("content")
  if hasError, message := services.ValidateModel(&article, r.PostForm); hasError == true {
    c.SetFlash(message)
  } else {
    c.DB.Save(&article)
    id = int(article.ID) // Because: https://github.com/jinzhu/gorm/blob/master/main.go#L390
  }
  http.Redirect(w, r, "/article/" + strconv.Itoa(id) + "/edit", 302)
}

/**
 * Remove Article (POST)
 */
func ArticleRemovePost(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var article models.Article
  c.DB.Find(&article, id)

  if article.ID < 1 || article.UserID != c.Page.User.ID {
    c.SetFlash("Not allowed to remove")
  } else {
    if err := c.DB.Unscoped().Delete(&article).Error; err != nil {
      c.SetFlash("Article was not removed")
    }
  }
  http.Redirect(w, r, "/articles", 302)
}

/**
 * View Profile
 */
func ProfileView(w http.ResponseWriter, r *http.Request, c *services.Containers) {
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