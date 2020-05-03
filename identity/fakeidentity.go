package identity

import (
  "net/http"
)

type FakeIdentity struct {
  userID string
}

func (id *FakeIdentity) GetUserID(r *http.Request) string {
  return id.userID
}
