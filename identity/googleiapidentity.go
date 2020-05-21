package identity

import (
  "fmt"
  "net/http"
)

type GoogleIAPIdentity struct {
}

func (id *GoogleIAPIdentity) GetUserID(r *http.Request) string {
  fmt.Printf("Request for identity: %+v\n", r)
  // TODO: This is not bullet proof and we should implement more defenses.
  // See https://cloud.google.com/go/getting-started/authenticate-users-with-iap
  return r.Header.Get("X-Goog-Authenticated-User-ID")
}
