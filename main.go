package main

import (
  "context"
  "errors"
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"

  "cloud.google.com/go/datastore"
)

const PROJECT_ID string = "pricebook"

type Item struct {
  Name string
  MinPrice float32
  Unit string
  // TODO: Which format for the date.
  // TODO: Add the timeseries.
}

type IDataStore interface {
  Add(it Item) (string, error)
  Delete(key string) error
  Update(key string, it Item) error
}

type GoogleDataStore struct {
  client *datastore.Client
}

func NewGoogleDataStore() *GoogleDataStore {
  // TODO: Singleton?
  ctx := context.Background()

  // Creates a client.
  client, err := datastore.NewClient(ctx, PROJECT_ID)
  if err != nil {
    log.Fatalf("Failed to create client: %v", err)
    return nil
  }
  return &GoogleDataStore{client}
}

func (ds *GoogleDataStore) Add(it Item) (string, error) {
  return "", errors.New("Not implemented")
}

func (ds *GoogleDataStore) Delete(key string) error {
  return errors.New("Not implemented")
}

func (ds *GoogleDataStore) Update(key string, it Item) error {
  return errors.New("Not implemented")
}


type InMemoryDataStore struct {
}

func NewInMemoryDataStore() *InMemoryDataStore {
  return &InMemoryDataStore{};
}

func (ds *InMemoryDataStore) Add(it Item) (string, error) {
  return "", errors.New("Not implemented")
}

func (ds *InMemoryDataStore) Delete(key string) error {
  return errors.New("Not implemented")
}

func (ds *InMemoryDataStore) Update(key string, it Item) error {
  return errors.New("Not implemented")
}


type defaultHandler struct {
  ds *IDataStore
}

func (defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  fmt.Fprintf(w, "Hello World!")
}

func main() {
  var localPtr = flag.Bool("local", false, "Set the binary for local testing")
  flag.Parse()

  var ds IDataStore
  if *localPtr {
    ds = NewInMemoryDataStore()
  } else {
    ds = NewGoogleDataStore()
  }

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("Defaulting to port %s", port)
  }

  http.Handle("/", defaultHandler{&ds})

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    log.Fatal(err)
  }
}


