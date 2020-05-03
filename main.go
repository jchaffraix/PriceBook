package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io"
  "log"
  "net/http"
  "os"

  "github.com/jchaffraix/PriceBook/datastore"
  "github.com/jchaffraix/PriceBook/identity"
)

type defaultHandler struct {
}

func (defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  fmt.Fprintf(w, "Hello World!")
}

type addHandler struct {
  ds datastore.IDataStore
}

func (h addHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // TODO: I should do user validation using
  // https://cloud.google.com/go/getting-started/authenticate-users-with-iap

  var it datastore.Item
  err := json.NewDecoder(r.Body).Decode(&it)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  key, err := h.ds.Add(it)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // TODO: Do I want a better format? Maybe json?
  fmt.Fprintf(w, key)
}

type deleteHandler struct {
  ds datastore.IDataStore
}

func (h deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // TODO: I should do user validation using
  // https://cloud.google.com/go/getting-started/authenticate-users-with-iap

  // The key is 64 bits, encoded in base 16 so 2 letters per byte.
  key := make([]byte, 128)
  n, err := r.Body.Read(key)

  // If there was any error, we just bail out. This is in contradiction
  // to the io.Reader documentation but it's the safe thing to do here.
  if err != nil && err != io.EOF {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = h.ds.Delete(string(key[:n]))
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  fmt.Fprintf(w, "Deleted")
}

type updateHandler struct {
  ds datastore.IDataStore
}

type UpdatePayload struct {
  ID string `json:"id"`
  It datastore.Item `json:"item"`
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // TODO: I should do user validation using
  // https://cloud.google.com/go/getting-started/authenticate-users-with-iap
  var payload UpdatePayload
  err := json.NewDecoder(r.Body).Decode(&payload)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  err = h.ds.Update(payload.ID, payload.It)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  fmt.Fprintf(w, "Updated")
}

func main() {
  var inMemoryPtr = flag.Bool("inmemory", false, "Toggle the datastore to in-memory for local testing")
  var fakeUserPtr = flag.String("fake_user_id", "", "Toggle FakeIdentity for local testing")
  flag.Parse()
  datastore.Init(*inMemoryPtr)
  identity.Init(*fakeUserPtr)

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("Defaulting to port %s", port)
  }

  http.Handle("/", defaultHandler{})
  http.Handle("/add", addHandler{datastore.Get()})
  http.Handle("/delete", deleteHandler{datastore.Get()})
  http.Handle("/update", updateHandler{datastore.Get()})

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    log.Fatal(err)
  }
}
