package main

import (
  "context"
  "fmt"
  "log"
  "net/http"
  "os"

  "cloud.google.com/go/datastore"
)

type defaultHandler struct {
  client *datastore.Client
}

func (defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  fmt.Fprintf(w, "Hello World!")
}

func main() {
  ctx := context.Background()

  // Set your Google Cloud Platform project ID.
  projectID := "pricebook"

  // Creates a client.
  client, err := datastore.NewClient(ctx, projectID)
  if err != nil {
    log.Fatalf("Failed to create client: %v", err)
  }

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("Defaulting to port %s", port)
  }

  http.Handle("/", defaultHandler{client})

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    log.Fatal(err)
  }
}


