package core

import (
  "github.com/jinzhu/gorm"
  "github.com/stanisdev/models"
  "html/template"
  "time"
)

type Containers struct {
  DB *gorm.DB
  Models models.StaticMethods
  Session SessionManager
  Page Page
  Params map[string]string
  BadRequest string
}

type Router struct {
  Handlers map[string]map[string]RouterHandler
  Config *Config
}

type Cookie struct {
  Name       string
  Value      string
  Path       string
  Domain     string
  Expires    time.Time
  RawExpires string
  MaxAge     int
  Secure     bool
  HttpOnly   bool
  Raw        string
  Unparsed   []string // Raw text of unparsed attribute-value pairs
}

type Page struct {
  Title string
  Flash template.HTML
  Url string
  User models.User
  Data map[string]interface{}
}

type Config struct {
  DbName string `json:"db_name"`
  DbUser string `json:"db_user"`
  DbPass string `json:"db_pass"`
  UrlsWithoutTemplate []string `json:"urls_without_template"`
  ProtectedUrls []struct{
    Url string 
    Method string
  } `json:"protected_urls"`
}

type paginationData struct {
  Value template.HTML
  Current bool
  Enabled bool
  Link int 
}

type typedRequestParams struct {
  Name string
  Type string
  DefaultValue interface{}
}