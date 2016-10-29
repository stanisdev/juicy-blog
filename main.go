package main

import (
  "github.com/stanisdev/core"
  "github.com/stanisdev/db"
  "github.com/stanisdev/helpers"
  "net/http"
  "fmt"
)

func init()  {
  db.Connect()
  core.SessionStart()
}

func main() {
  c := &helpers.Containers{DB: db.GetConnection()}
  http.HandleFunc("/", helpers.MakeHandler(core.Index, c))
  http.HandleFunc("/article/add", helpers.MakeHandler(core.AddArticle, c))

  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}
