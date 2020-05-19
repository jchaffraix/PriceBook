package datastore

import (
  "log"
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
  err := ValidateItemWithoutKey(it)
  if err != nil {
    return "", err
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

func (ds *InMemoryDataStore) Update(userID string, item Item) error {
  err := ValidateItemWithKey(item)
  if err != nil {
    return err
  }

  fullKey := buildKey(userID, item.ID)
  _, found := ds.m[fullKey]
  if found {
    // Remove the ID field to match the behavior of the GoogleDataStore.
    item.ID = ""
    ds.m[fullKey] = item
    return nil
  }
  return &NotFoundError{fullKey}
}

func (ds *InMemoryDataStore) Get(userID string) []Item {
  res := make([]Item, 0)
  for key, it := range (ds.m) {
    keyUserID, keyID := splitKey(key)
    if keyUserID == userID {
      err := ValidateItemWithoutKey(it)
      if err != nil {
        log.Printf("Item in storage doens't pass validation: %+v", it)
        continue
      }
      it.ID = keyID
      res = append(res, it)
    }
  }
  return res
}
