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

  http.ServeFile(w, r, "index.html")
}

type addHandler struct {
  ds datastore.IDataStore
  id identity.IIdentity
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
  userID := h.id.GetUserID(r)
  key, err := h.ds.Add(userID, it)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // TODO: Do I want a better format? Maybe json?
  fmt.Fprintf(w, key)
}

type deleteHandler struct {
  ds datastore.IDataStore
  id identity.IIdentity
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

  userID := h.id.GetUserID(r)
  err = h.ds.Delete(userID, string(key[:n]))
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  fmt.Fprintf(w, "Deleted")
}

type updateHandler struct {
  ds datastore.IDataStore
  id identity.IIdentity
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // TODO: I should do user validation using
  // https://cloud.google.com/go/getting-started/authenticate-users-with-iap
  var item datastore.Item
  err := json.NewDecoder(r.Body).Decode(&item)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  userID := h.id.GetUserID(r)
  log.Printf("Preparing to update object %+v for userID %v", item, userID)
  err = h.ds.Update(userID, item)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  fmt.Fprintf(w, "Updated")
}

type getHandler struct {
  ds datastore.IDataStore
  id identity.IIdentity
}

func (h getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  userID := h.id.GetUserID(r)
  items := h.ds.Get(userID)
  res, err := json.Marshal(items)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Add("Content-Type", "application/json")
  w.Write(res)
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
  http.Handle("/add", addHandler{datastore.Get(), identity.Get()})
  http.Handle("/delete", deleteHandler{datastore.Get(), identity.Get()})
  http.Handle("/update", updateHandler{datastore.Get(), identity.Get()})
  http.Handle("/get", getHandler{datastore.Get(), identity.Get()})

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    log.Fatal(err)
  }
}
