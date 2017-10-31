package services

import (
  redis "gopkg.in/redis.v5"
  "net/http"
  "time"
)

type SessionManager struct {
  client *redis.Client
  sid string
}

func (sm *SessionManager) Start(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("sid")
  if len(cookie.String()) < 1 { // Set sid
    sm.createCookie(w)
  } else {
    sm.sid = cookie.Value
  }
  sm.client = redis.NewClient(&redis.Options{
    Addr: "127.0.0.1:6379",
    Password: "",
    DB: 0,
  })
}

func (sm *SessionManager) createCookie(w http.ResponseWriter) {
  sid := GenerateRandomString(40)
  expiration := time.Now().Add(365 * 24 * time.Hour)
  cookie := http.Cookie{
    Name: "sid",
    Value: sid,
    Expires: expiration,
  }
  http.SetCookie(w, &cookie)
  sm.sid = sid
}

func (sm *SessionManager) Set(key string, value string) {
  sm.client.HSet(sm.sid, key, value)
}

func (sm *SessionManager) Get(key string) (string, bool) {
  res, err := sm.client.HGet(sm.sid, key).Result()
  if err != nil {
    return "", false
  } else {
    return res, len(res) > 0
  }
}

func (sm *SessionManager) Unset(key string) {
  sm.client.HDel(sm.sid, key)
}
