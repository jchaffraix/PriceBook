package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"

  "github.com/jchaffraix/PriceBook/datastore"
)

type defaultHandler struct {
  ds *datastore.IDataStore
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

  var ds datastore.IDataStore
  if *localPtr {
    ds = datastore.NewInMemoryDataStore()
  } else {
    ds = datastore.NewGoogleDataStore()
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


