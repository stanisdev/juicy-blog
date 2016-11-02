package main

import (
  "github.com/stanisdev/core"
  "net/http"
  "fmt"
)

func main() {
  core.DatabaseMigrate()
  http.HandleFunc("/", core.MakeHandler(core.Index))
  http.HandleFunc("/login", core.MakeHandler(core.Login))
  http.HandleFunc("/article/add", core.MakeHandler(core.AddArticle))  

  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}
