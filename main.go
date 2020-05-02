package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("Defaulting to port %s", port)
  }

  http.HandleFunc("/", defaultHandler)

  if err := http.ListenAndServe(":"+port, nil); err != nil {
    log.Fatal(err)
  }
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  fmt.Fprintf(w, "Hello World!")
}
