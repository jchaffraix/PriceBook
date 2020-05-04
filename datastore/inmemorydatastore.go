package datastore

import (
  "errors"
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

func buildKey(userID, key string) string {
  return userID + "." + key
}

func splitKey(fullKey string) (string, string) {
  split := strings.Split(fullKey, ".")
  if len(split) != 2 {
    return "", ""
  }
  return split[0], split[1]
}

func (ds *InMemoryDataStore) Add(userID string, it Item) (string, error) {
  if it.ID != "" {
    return "", &InvalidItemError{"Unexpected ID in Add"}
  }

  if it.Name == "" {
    return "", &InvalidItemError{"Missing 'name'"}
  }

  key := randomKey()
  fullKey := buildKey(userID, key)
  ds.m[fullKey] = it
  return key, nil
}

func (ds *InMemoryDataStore) Delete(userID, key string) error {
  fullKey := buildKey(userID, key)
  _, found := ds.m[fullKey]
  if found {
    delete(ds.m, fullKey)
    return nil
  }
  return &NotFoundError{fullKey}
}

func (ds *InMemoryDataStore) Update(userID string, it Item) error {
  if it.ID == "" {
    return errors.New("Missing key for update")
  }

  fullKey := buildKey(userID, it.ID)
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
    keyUserID, keyID := splitKey(key)
    if keyUserID == userID {
      it.ID = keyID
      res = append(res, it)
    }
  }
  return res
}
