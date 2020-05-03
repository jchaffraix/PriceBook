package datastore

import (
  "math/rand"
  "strings"
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

func (ds *InMemoryDataStore) Add(userID string, it Item) (string, error) {
  if it.Name == "" {
    return "", &InvalidItemError{"Missing 'name'"}
  }

  key := randomKey()
  fullKey := userID + key
  ds.m[fullKey] = it
  return key, nil
}

func (ds *InMemoryDataStore) Delete(userID, key string) error {
  fullKey := userID + key
  _, found := ds.m[fullKey]
  if found {
    delete(ds.m, fullKey)
    return nil
  }
  return &NotFoundError{fullKey}
}

func (ds *InMemoryDataStore) Update(userID, key string, it Item) error {
  fullKey := userID + key
  _, found := ds.m[fullKey]
  if found {
    ds.m[fullKey] = it
    return nil
  }
  return &NotFoundError{fullKey}
}

func (ds *InMemoryDataStore) Get(userID string) []Item {
  var res []Item
  for key, it := range (ds.m) {
    if strings.HasPrefix(key, userID) {
      res = append(res, it)
    }
  }
  return res
}
