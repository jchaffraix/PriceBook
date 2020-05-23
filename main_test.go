package main

import (
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "testing"
  "strings"

  "github.com/jchaffraix/PriceBook/datastore"
  "github.com/jchaffraix/PriceBook/identity"
)

func initModules() {
  datastore.Init(true)
  identity.Init("1234")
}

func TestAddFailWithNoBody(t *testing.T) {
  initModules()

  req := httptest.NewRequest("GET", "http://example.com/add", nil)
  rec := httptest.NewRecorder()

  handler := addHandler{datastore.Get(), identity.Get()}
  handler.ServeHTTP(rec, req)

  res := rec.Result()
  if res.StatusCode != http.StatusBadRequest {
    t.Errorf("Error code is wrong (expected 400): %+v", res.StatusCode)
  }
  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)
  bodyStr := strings.TrimSpace(string(body))
  const expectedStr string = "Bad JSON"
  if bodyStr != expectedStr {
    t.Errorf("Unexpected body. Expected %v, got %v", expectedStr, bodyStr)
  }
}

func TestAddAllowValidObject(t *testing.T) {
  initModules()

  tt := []struct {
    name string
    body string
  }{
    {"Simple object", "{\"name\":\"Oranges\",\"brand\":\"Navel\",\"quantity\":\"1\",\"unit\":\"lb\",\"purchases\":[{\"date\":\"2020-05-1\",\"store\":\"Foobar\",\"price\":\"1.29\",\"currency\":\"$\"}]}"},
  }
  for _, tc := range tt {
    t.Run(tc.name, func(t *testing.T) {
      reqBody := strings.NewReader(tc.body)
      req := httptest.NewRequest("GET", "http://example.com/add", reqBody)
      rec := httptest.NewRecorder()

      handler := addHandler{datastore.Get(), identity.Get()}
      handler.ServeHTTP(rec, req)

      res := rec.Result()
      if res.StatusCode != http.StatusOK {
        t.Errorf("Error code is wrong (expected 200): %+v", res.StatusCode)
      }
      defer res.Body.Close()
      body, _ := ioutil.ReadAll(res.Body)
      bodyStr := strings.TrimSpace(string(body))
      if bodyStr == "" {
        t.Errorf("Expected body. got empty response")
      }
      // TODO: Get it back.
    })
  }
}
