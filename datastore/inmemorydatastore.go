package datastore

import (
  "math/rand"
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
  _, found := ds.m[key]
  if found {
    delete(ds.m, key)
    return nil
  }
  return &NotFoundError{key}
}

func (ds *InMemoryDataStore) Update(key string, it Item) error {
  _, found := ds.m[key]
  if found {
    ds.m[key] = it
    return nil
  }
  return &NotFoundError{key}
}



