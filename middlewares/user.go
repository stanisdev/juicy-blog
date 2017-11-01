package middlewares

import (
  "github.com/stanisdev/juicy-blog/services"
)

func Auth(c *services.Containers) {
  if c.Page.User.Authorized() {
    c.Next()
  }
}
