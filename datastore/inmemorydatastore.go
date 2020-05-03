package datastore

import (
  "errors"
)


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



