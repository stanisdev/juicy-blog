package main

import (
  "github.com/stanisdev/core"
  "github.com/stanisdev/db"
  "net/http"
  "fmt"
)

func init()  {
  db.Connect()
}

func main() {
  dbConnection := db.GetConnection()
  session := core.SessionManager{DB: dbConnection}
  c := &core.Containers{DB: dbConnection, Session: session}
  http.HandleFunc("/", core.MakeHandler(core.Index, c))
  http.HandleFunc("/article/add", core.MakeHandler(core.AddArticle, c))

  fmt.Println("Server is listening on port: 8080")
  http.ListenAndServe(":8080", nil)
}
