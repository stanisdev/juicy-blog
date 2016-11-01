package main

import (
  "github.com/stanisdev/core"
  "net/http"
  "fmt"
)

func main() {
  http.HandleFunc("/", core.MakeHandler(core.Index))
  http.HandleFunc("/article/add", core.MakeHandler(core.AddArticle))

  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}
