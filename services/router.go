package services

import (
  "fmt"
  "net/http"
  "regexp"
  "strings"
  "runtime"
  "path"
  "io/ioutil"
  "encoding/json"
  "github.com/stanisdev/juicy-blog/models"
  "html/template"
  "os"
  "reflect"
  "time"
)

type Router struct {
  Handlers map[string] map[string] RouterData
  Config *Config
  ProjectDir string
}

type RouterHandler func(http.ResponseWriter, *http.Request, *Containers)

type RouterData struct {
  Handler RouterHandler
  Middlewares []func(*Containers)
}

type Flash struct {
  Message template.HTML
  State string
}

type Page struct {
  Title string
  Flash Flash
  Url string
  User models.User
  Data map[string]interface{}
}

type Config struct {
  DbName string `json:"db_name"`
  DbUser string `json:"db_user"`
  DbPass string `json:"db_pass"`
  CommentsByPage int `json:"comments_by_page"`
}

/**
 * Init router by loading and preparing basic data
 */
func (r *Router) Init() {
  _, filename, _, ok := runtime.Caller(0)
  if !ok {
    panic("No caller information")
  }
  r.ProjectDir = path.Dir(path.Dir(filename))
  fmt.Println(r.ProjectDir)
  // Load Config
  raw, err := ioutil.ReadFile(r.ProjectDir + "/config.json")
  if err != nil {
    panic("Config file cannot be loaded")
  }
  var config Config
  if err := json.Unmarshal(raw, &config); err != nil {
    panic("JSON config cannot be parsed")
  }
  r.Config = &config

  // Public file server
  http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(r.ProjectDir + "/public"))))

  if len(os.Getenv("DB_MIGRATE")) > 0 {
    DatabaseMigrate(r.Config)
  }
  if len(os.Getenv("IMPORT_DEMO_FIXTURES")) > 0 {
    ImportDemoFixtures(r.Config)
  }
}

/**
 * Auxiliary method
 */
func (r *Router) defineHandler(url string) {
  if len(r.Handlers[url]) < 1 {
    r.Handlers[url] = make(map[string]RouterData)
  }
}

/**
 * Adding GET handler
 */
func (r *Router) GET(url string, fn RouterHandler, middlewares ...func(*Containers)) {
  r.defineHandler(url)
  routerData := RouterData{
    Handler: fn,
  }
  if len(middlewares) > 0 {
    routerData.Middlewares = middlewares
  }
  r.Handlers[url]["GET"] = routerData
}

/**
 * Adding POST handler
 */
func (r *Router) POST(url string, fn RouterHandler, middlewares ...func(*Containers)) {
  r.defineHandler(url)
  routerData := RouterData{
    Handler: fn,
  }
  if len(middlewares) > 0 {
    routerData.Middlewares = middlewares
  }
  r.Handlers[url]["POST"] = routerData
}

/**
 * Page not found handle
 */
func (r *Router) notFound(w http.ResponseWriter) {
  w.WriteHeader(http.StatusNotFound)
  fmt.Fprint(w, "Page not found 404")
}

/**
 * Primary handler
 */
func (self *Router) handler(w http.ResponseWriter, r *http.Request) {
  url := r.URL.Path // Current URL
  if url == "/favicon.ico" {
    return
  }
  reqData := make(map[string]string)
  handlersByUrlPattern, isUrlExists := self.Handlers[url]
  if !isUrlExists { // First compare by simple equal
    var urlParamsValues []string
    reg, _ := regexp.Compile("\\/:([a-z]+\\/??)")
    for urlPattern, _ := range self.Handlers { // Compare to patterns
      if strings.Contains(urlPattern, ":") {

        preparedUrlPattern := reg.ReplaceAllString(urlPattern, "(\\/[^\\/]+)")
        re, _ := regexp.Compile("^" + preparedUrlPattern + "$")
        urlParamsValues = re.FindStringSubmatch(url)

        if len(urlParamsValues) > 1 { // Pattern was found
          urlParamsValues = urlParamsValues[1:]
          handlersByUrlPattern, _ = self.Handlers[urlPattern]
          urlNamedParams := reg.FindAllString(urlPattern, len(urlParamsValues) )

          for key, val := range urlNamedParams { // Loop through url named params
            var reqValue string = urlParamsValues[key]
            if len(reqValue) > 0 {
              reqValue = reqValue[1:]
            }
            reqData[val[2:]] = reqValue
          }
          break
        }
      }
    }
    if len(urlParamsValues) < 1 {
      self.notFound(w)
      return
    }
  }
  // Get necessary handler
  routerData, isHandlerExists := handlersByUrlPattern[r.Method]
  if !isHandlerExists {
    self.notFound(w)
    return
  }
  handler := &routerData.Handler

  // Connect to DB
  dbConnection := ConnectToDatabase(self.Config)

  // Prepare data container
  container := &Containers {
    UrlParams: reqData,
    DB: dbConnection,
    Models: models.DatabaseStaticMethods{ DB: dbConnection },
    Page: Page{ Data: make(map[string] interface{}) },
    Session: SessionManager{},
    ResponseWriter: &w,
    Request: r,
    Config: self.Config,
    IntUrlParams: make(map[string]int),
    MiddlewaresData: make(map[string]interface{}),
  }
  container.Session.Start(w, r)
  container.Auth()

  defer func() {
    if message := recover(); message != nil {
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprint(w, message)
    }
  }()

  // Middlewares
  middlewares := routerData.Middlewares
  if len(middlewares) > 0 {
    for _, middleware := range middlewares { // Run every middleware
      container.UrlIsRestricted = true
      middleware(container)
    }
    if container.UrlIsRestricted { // Check successfully
      self.notFound(w)
      return
    }
  }

  (*handler)(w, r, container) // Run url-handler

  // Some intolerable errors
  if len(container.BadRequest) > 0 {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w, container.BadRequest)
    return
  }

  // Render template
  if r.Method == "GET" {
    if container.NoTemplate {
      return
    }
    message, state, hasFlash := container.GetFlash()
    if (hasFlash == true) {
      container.Page.Flash = Flash {
        Message: template.HTML(message),
        State: state,
      }
    }
    container.Page.Url = url

    // Prepare info for template
    var templateName string
    if len(container.SpecifiedTemplate) < 1 {
      handlerName := runtime.FuncForPC(reflect.ValueOf(*handler).Pointer()).Name()
      templateName = handlerName[strings.LastIndex(handlerName, ".") + 1:]
    } else {
      templateName = container.SpecifiedTemplate
    }
    self.loadTemplate(strings.ToLower(templateName), w, &container.Page)
  }
}

/**
 * Starting app
 */
func (r *Router) loadTemplate(templateName string, w http.ResponseWriter, p *Page) {
  w.Header().Set("Content-type", "text/html")
  tplFuncMap := make(template.FuncMap)
  tplFuncMap["IsArticles"] = func (url string) bool {
    return len(url) > 7 && url[:8] == "/article"
  }
  tplFuncMap["DateFormat"] = func (date time.Time) string {
    return date.Format("_2 Jan 2006 15:04:05")
  }
  tplFuncMap["ChangeToBr"] = func (value string) interface{} {
    return template.HTML(strings.Replace(value, "\n", "<br/>", -1))
  }
  layoutPath := r.ProjectDir + "/templates/layouts/base.html"
  templatePath := r.ProjectDir + "/templates/" + templateName + ".html"
  t, err := template.New("").Funcs(tplFuncMap).ParseFiles(layoutPath, templatePath)
  if err != nil {
    fmt.Fprintf(w, "Template cannot be loaded")
  }
  t.ExecuteTemplate(w, "base", p)
}

/**
 * Starting app
 */
func (r *Router) Run() {
  http.HandleFunc("/", r.handler)
  fmt.Println("---------------- started ----------------")
  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}
