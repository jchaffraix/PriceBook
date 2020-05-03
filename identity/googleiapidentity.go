package identity

import (
  "net/http"
)

type GoogleIAPIdentity struct {
}

func (id *GoogleIAPIdentity) GetUserID(r *http.Request) string {
  // TODO: This is not bullet proof and we should implement more defenses.
  // See https://cloud.google.com/go/getting-started/authenticate-users-with-iap
  return r.Header.Get("X-Goog-Authenticated-User-ID")
}
