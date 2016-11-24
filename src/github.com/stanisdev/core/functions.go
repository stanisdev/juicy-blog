package core

import (
  "net/http"
  "time"
  "math/rand"
  "html/template"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/url"
  "github.com/gorilla/schema"
  validator "github.com/asaskevich/govalidator"
  "strings"
  "strconv"
)

const viewPath = "src/github.com/stanisdev/templates/"

/**
 * Container's methods
 */
func (c *Containers) SetFlash(value string) {
  c.Session.Set("flash", value)
}

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

func (c *Containers) GetParamByType(data typedRequestParam) (success bool, result interface{}) {
  var name = data.Name
  switch data.Type {
    case "string":
      success = len(c.Params[name]) > 0
      result = c.Params[name]
    case "int":
      if len(c.Params[name]) > 0 {
        value, err := strconv.Atoi(c.Params[name])
        if err != nil {
          c.BadRequest = strings.Title(name) + " parameter must be a number"
        } else if value < 1 {
          c.BadRequest = strings.Title(name) + " parameter must be greater then zero"
        } else {
          success = true
          result = value
        }
      } else {
        success = true
        result = data.DefaultValue
      }
  }
  return
}

/**
 * General Functions
 */
func GetConfig() *Config {
  raw, err := ioutil.ReadFile("./config.json")
  if err != nil {
    panic("Config file cannot be loaded")
  }
  var config Config
  if err := json.Unmarshal(raw, &config); err != nil {
    panic("JSON config cannot be parsed")
  }
  return &config
}

func loadTemplate(templateName string, w http.ResponseWriter, p *Page)  {
  w.Header().Set("Content-type", "text/html")
  tplFuncMap := make(template.FuncMap)
  tplFuncMap["IsArticles"] = func (url string) bool {
    return len(url) > 8 && url[:9] == "/articles"
  }
  tplFuncMap["DateFormat"] = func (date time.Time) string {
    return date.Format("_2 Jan 2006 15:04:05")
  }
  tplFuncMap["ChangeToBr"] = func (value string) interface{} {
    return template.HTML(strings.Replace(value, "\n", "<br/>", -1))
  }
  t, err := template.New("").Funcs(tplFuncMap).ParseFiles(viewPath + "/layouts/layout.html", viewPath + templateName + ".html")
  if err != nil {
    fmt.Fprintf(w, "Template cannot be loaded")
  }
  t.ExecuteTemplate(w, "layout", p)
}

func GenerateRandomString(l int) string {
  rand.Seed(time.Now().UTC().UnixNano())
  bytes := make([]byte, l)
  for i := 0; i < l; i++ {
      bytes[i] = byte(randInt(65, 90))
  }
  return string(bytes)
}

func randInt(min int, max int) int {
  return min + rand.Intn(max-min)
}

func ValidateModel(modelInstance interface{}, formData url.Values) (bool, string) {
  decoder := schema.NewDecoder()
  if err := decoder.Decode(modelInstance, formData); err != nil {
    return true, "Fields cannot be parsed"
  }
  if _, err := validator.ValidateStruct(modelInstance); err != nil {
    var message string = err.Error()
    var splited []string = strings.Split(message[:len(message)-1], ";")
    return true, strings.Join(splited, "<br/>")
  } else {
    return false, ""
  }
}

func makePagination(curPage int, pagesCount int) []paginationData {
  middle := make([]int, 0, 4)
  if curPage == 1 { // First
    middle = append(middle, 1, 2)
  } else if curPage == pagesCount { // Last
    middle = append(middle, pagesCount - 1, pagesCount)
  } else { // Intermediate
    middle = append(middle, curPage - 1, curPage, curPage + 1)
    if curPage == 3 {
      middle = append([]int{1}, middle...)
    } else if curPage == (curPage - 2) {
      middle = append(middle, pagesCount)
    }
  }
  prepare := make([]int, 1, 10)
  prepare[0] = -1
  if middle[0] > 1 { // Add first page
    prepare = append(prepare, 1)
    if pagesCount > 3 {
      prepare = append(prepare, 0)
    }
  }
  prepare = append(prepare, middle...)
  if middle[len(middle) - 1] < pagesCount { // Add last page
    if pagesCount > 3 {
      prepare = append(prepare, 0)
    }
    prepare = append(prepare, pagesCount)
  }
  prepare = append(prepare, -2)
  result := make([]paginationData, len(prepare))
  for key, value := range prepare {
    switch value {
    case -1:
        result[key].Value = template.HTML("&laquo;")
        if curPage != 1 {
          result[key].Link = curPage - 1
          result[key].Enabled = true
        }
    case 0:
        result[key].Value = "..." 
    case -2:
        result[key].Value = template.HTML("&raquo;")
        if curPage < pagesCount {
          result[key].Link = curPage + 1
          result[key].Enabled = true
        }
    default:
        result[key].Value = template.HTML(strconv.Itoa(value))
        if value == curPage {
          result[key].Current = true
        } else {
          result[key].Link = value
          result[key].Enabled = true
        }
    }
  }
  return result
}