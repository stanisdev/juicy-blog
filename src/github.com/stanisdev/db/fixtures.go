package db

import (
  "github.com/stanisdev/models"
  "github.com/jinzhu/gorm"
  "fmt"
)

func ImportFixtures(config ...string) {
  var con *gorm.DB = Connect(config[0], config[1], config[2]);
  fixtures := getFixtures()
  ch := make(chan bool)
  for _, user := range fixtures["Users"] {
    go func(user map[string]interface{}) {
      var newUser models.User
      newUser.ID = uint(user["id"].(int))
      newUser.Name = user["name"].(string)
      newUser.Email = user["email"].(string)
      newUser.Password = user["password"].(string)
      con.Create(&newUser)
      ch <- true
    }(user)
  }
  for i := 0; i < 2; i++ {
    <-ch
  }
  close(ch)
  fmt.Println("All fixtures have been imported successfully")
}

func getFixtures() map[string][]map[string]interface{} {
  return map[string][]map[string]interface{} {
    "Users": []map[string]interface{} {
      {
        "id": 1,
        "name": "Stan",
        "email": "stan@gmail.com", 
        "password": "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
      },
      {
        "id": 2, 
        "name": "Abrasha",
        "email": "abrasha@gmail.com", 
        "password": "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
      },
    },
  }
}