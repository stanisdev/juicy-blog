package middlewares

import (
  "github.com/stanisdev/juicy-blog/services"
)

func Auth(c *services.Containers) {
  if c.User.Authorized() {
    c.Next()
  }
}
