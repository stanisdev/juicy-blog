package services

import (
  "github.com/jinzhu/gorm"
  "github.com/stanisdev/juicy-blog/models"
  "strconv"
  "strings"
)

type TypedRequestParam struct {
  Name string
  Type string
  DefaultValue interface{}
}

type Containers struct {
  UrlParams map[string]string
  DB *gorm.DB
  Models models.DatabaseStaticMethods
  Page Page
  Session SessionManager
  BadRequest string
  UrlIsRestricted bool
}

/**
 * Set flash message
 */
func (c *Containers) SetFlash(value string) {
  c.Session.Set("flash", value)
}

/**
 * Get flash message
 */
func (c *Containers) GetFlash() (string, bool) {
  value, exists := c.Session.Get("flash")
  if exists == true {
    c.Session.Unset("flash")
    return value, true
  } else {
    return "", false
  }
}

func (c *Containers) Auth() {
  userId, isAuth := c.Session.Get("user")
  if isAuth == true {
    c.DB.Select("id, name, email").First(&c.Page.User, userId)
  }
}

func (c *Containers) GetParamByType(data TypedRequestParam) (result interface{}) {
  var name = data.Name
  switch data.Type {
    case "string":
      result = c.UrlParams[name]
    case "int":
      if len(c.UrlParams[name]) > 0 {
        value, err := strconv.Atoi(c.UrlParams[name])
        if err != nil {
          panic(strings.Title(name) + " parameter must be a number")
        } else if value < 1 {
          panic(strings.Title(name) + " parameter must be greater then zero")
        } else {
          result = value
        }
      } else {
        result = data.DefaultValue
      }
  }
  return
}

/**
 * To confirm middleware successfully passing checking
 */
func (c *Containers) Next() {
  c.UrlIsRestricted = false
}
