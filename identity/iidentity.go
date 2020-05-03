package identity

import (
  "net/http"
)

type IIdentity interface {
  GetUserID(r *http.Request) string
}

var id IIdentity
func Init(userID string) {
  if userID != "" {
    id = &FakeIdentity{userID}
    return
  }
  id = &GoogleIAPIdentity{}
}

func Get() IIdentity {
  return id
}
