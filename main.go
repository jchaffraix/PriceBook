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
}

func (defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  fmt.Fprintf(w, "Hello World!")
}

func main() {
  var inMemoryPtr = flag.Bool("inmemory", false, "Toggle the datastore to in-memory for local testing")
  flag.Parse()
  datastore.Init(*inMemoryPtr);

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("Defaulting to port %s", port)
  }

  http.Handle("/", defaultHandler{})

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    log.Fatal(err)
  }
}
