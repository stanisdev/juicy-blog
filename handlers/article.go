package handlers

import (
  "net/http"
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/models"
  "sync"
  "math"
  "strconv"
  "time"
  "fmt"
)

/**
 * Index page
 */
func Index(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  var articlesCount, usersCount, newArticlesCount int
  waitGroup := sync.WaitGroup{}
  waitGroup.Add(3)
  go func() {
    c.DB.Model(models.Article{}).Count(&articlesCount)
    waitGroup.Done()
  }()
  go func() {
    c.DB.Model(models.User{}).Count(&usersCount)
    waitGroup.Done()
  }()
  go func() {
    c.DB.Model(models.NewArticlesSubscriber{}).Where("subscriber_id = ?", c.User.ID).Count(&newArticlesCount)
    waitGroup.Done()
  }()
  waitGroup.Wait()
  c.Page.Data["ArticlesCount"] = articlesCount
  c.Page.Data["UsersCount"] = usersCount
  c.Page.Data["NewArticlesCount"] = newArticlesCount
  c.Page.Title = "Mutual Blog"
}


func Articles(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  page := c.GetParamByType(services.TypedRequestParam{ Name: "page", Type: "int", DefaultValue: 1 }).(int)
  if page < 1 {
    c.BadRequest = "Page must be greater then zero"
    return
  }
  filter := r.URL.Query().Get("filter")
  c.Page.Data["filter"] = filter
  if len(filter) < 1 {
    filter = "common"
  }
  var offset int = 0
  if page > 1 {
    offset = (page - 1) * 5
  }
  queryParams := map[string]int {
    "limit": 5,
    "offset": offset,
    "currUserId": int(c.User.ID),
  }
  data := c.Models.GetArticles(queryParams, filter)
  var count int = data["count"].(int)

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
  c.Page.Data["articles"] = data["articles"]
  c.Page.Data["filterCaption"] = data["caption"]
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
  article.UserID = c.User.ID
  if hasError, message := services.ValidateModel(&article, r.PostForm); hasError == true {
    c.SetFlash(message, "danger")
    c.Redirect("/articles/new")
  } else {
    if err := c.DB.Create(&article).Error; err != nil {
      c.SetFlash("Article cannot be created", "danger")
      c.Redirect("/articles/new")
    } else {
      var subscriberIds []uint = c.Models.FindAllSubscriberIds(c.User.ID)
      c.Models.AddNotifications(subscriberIds, article.ID)

      c.SetFlash("Article was created", "info")
      c.Redirect("/article/" + strconv.Itoa(int(article.ID)))
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
    c.Page.Data["canEdit"] = article.Userid == c.User.ID
  }
}

/**
 * Edit Article (GET)
 */
func ArticleEdit(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var article models.Article
  c.DB.Find(&article, id)

  if article.ID < 1 || article.UserID != c.User.ID {
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
  if article.ID < 1 || article.UserID != c.User.ID {
    c.SetFlash("Not allowed to edit", "danger")
    c.Redirect("/articles")
    return
  }
  r.ParseForm()
  article.Title = r.FormValue("title")
  article.Content = r.FormValue("content")
  if hasError, message := services.ValidateModel(&article, r.PostForm); hasError == true {
    c.SetFlash(message, "danger")
  } else {
    c.DB.Save(&article)
    id = int(article.ID) // Because: https://github.com/jinzhu/gorm/blob/master/main.go#L390
    c.SetFlash("Article was edited", "info")
  }
  c.Redirect("/article/" + strconv.Itoa(id))
}

/**
 * Remove Article (POST)
 */
func ArticleRemovePost(w http.ResponseWriter, r *http.Request, c *services.Containers) {
  id := c.GetParamByType(services.TypedRequestParam{ Name: "id", Type: "int", DefaultValue: nil }).(int)
  var article models.Article
  c.DB.Find(&article, id)

  if article.ID < 1 || article.UserID != c.User.ID {
    c.SetFlash("Not allowed to remove", "danger")
  } else {
    if err := c.DB.Unscoped().Delete(&article).Error; err != nil {
      c.SetFlash("Article was not removed", "danger")
    } else {
      c.SetFlash("Article was removed", "info")
    }
  }
  fmt.Println("Done")
  c.Redirect("/articles")
}
