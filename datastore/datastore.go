package datastore

import (
  "context"
  "errors"
  "log"
  "cloud.google.com/go/datastore"
)


const PROJECT_ID string = "pricebook"

type Item struct {
  Name string
  Quantity float32
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



