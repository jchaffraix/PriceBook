package datastore

import (
  "math/rand"
  "errors"
)


type InMemoryDataStore struct {
  m map[string] Item
}

func NewInMemoryDataStore() *InMemoryDataStore {
  return &InMemoryDataStore{make(map[string] Item)};
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// This generate a key for each insertion.
// Note that we intentionally do not set the seed for RNG
// as it makes testing more reliable.
func randomKey() string {
  b := make([]rune, 64)
  for i := range b {
    b[i] = letterRunes[rand.Intn(len(letterRunes))]
  }
  return string(b)
}

func (ds *InMemoryDataStore) Add(it Item) (string, error) {
  if it.Name == "" {
    return "", &InvalidItemError{"Missing 'name'"}
  }

  key := randomKey()
  ds.m[key] = it
  return key, nil
}

func (ds *InMemoryDataStore) Delete(key string) error {
  return errors.New("Not implemented")
}

func (ds *InMemoryDataStore) Update(key string, it Item) error {
  return errors.New("Not implemented")
}



