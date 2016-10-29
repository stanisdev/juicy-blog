package helpers

import (
  "github.com/jinzhu/gorm"
  "net/http"
)

type Containers struct {
  DB *gorm.DB
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *Containers), c *Containers) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request)  {
    fn(w, r, c)
  }
}
