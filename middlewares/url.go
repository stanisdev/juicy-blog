package middlewares

import (
  "github.com/stanisdev/juicy-blog/services"
  "strconv"
  "strings"
)

func ValidateUrl(params map[string]string) func(*services.Containers) {
  return func (c *services.Containers) {

    // Iterate throght params and check type conformity
    for paramName, paramType := range params {
      switch paramType {
        case "int":
          value, err := strconv.Atoi(c.UrlParams[paramName])
          if err != nil {
            panic(strings.Title(paramName) + " parameter must be a number")
          } else if value < 1 {
            panic(strings.Title(paramName) + " parameter must be greater then zero")
          } else {
            c.IntUrlParams[paramName] = value
          }
        default:
          panic("The " + paramType + " type is not processesed. You can add determine it in switch-condition")
      }
    }
    c.Next()
  }
}
