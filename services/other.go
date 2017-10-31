package services

import (
  "time"
  "math/rand"
  "html/template"
  "strconv"
  "github.com/gorilla/schema"
  validator "github.com/asaskevich/govalidator"
  "strings"
  "net/url"
)

type PaginationData struct {
  Value template.HTML
  Current bool
  Enabled bool
  Link int
}

/**
 * Blah-blah
 */
func GenerateRandomString(l int) string {
  rand.Seed(time.Now().UTC().UnixNano())
  bytes := make([]byte, l)
  for i := 0; i < l; i++ {
      bytes[i] = byte(randInt(65, 90))
  }
  return string(bytes)
}

/**
 * Blah-blah
 */
func randInt(min int, max int) int {
  return min + rand.Intn(max-min)
}

/**
 * Blah-blah
 */
func MakePagination(curPage int, pagesCount int) []PaginationData {
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
  result := make([]PaginationData, len(prepare))
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
