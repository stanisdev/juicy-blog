package middlewares

import (
  "github.com/stanisdev/juicy-blog/services"
  "github.com/stanisdev/juicy-blog/models"
)

func GetAritcle(c *services.Containers) {
  id, _ := c.IntUrlParams["id"]
  var article models.Article

  if c.DB.Find(&article, id).RecordNotFound() {
    c.SetFlash("Article not found", "danger")
    c.Redirect("/articles")
    return
  }
  c.MiddlewaresData["article"] = article
  c.Next()
}
