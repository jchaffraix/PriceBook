package datastore

import (
  "context"
  "errors"
  "log"

  "cloud.google.com/go/datastore"
)

const PROJECT_ID string = "pricebook"

type GoogleDataStore struct {
  client *datastore.Client
}

func NewGoogleDataStore() *GoogleDataStore {
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
