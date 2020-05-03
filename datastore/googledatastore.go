package datastore

import (
  "context"
  "log"
  "strconv"

  "cloud.google.com/go/datastore"
)

const PROJECT_ID string = "pricebook"
const TABLE string = "Items"

type GoogleDataStore struct {
  ctx context.Context
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
  return &GoogleDataStore{ctx, client}
}

func (ds *GoogleDataStore) Add(it Item) (string, error) {
  // TODO: Share the validation code with the InMemoryDataStore.
  if it.Name == "" {
    return "", &InvalidItemError{"Missing 'name'"}
  }

  // TODO: Pass the user instead of nil to scope their request.
  newKey := datastore.IncompleteKey(TABLE, nil)
  key, err := ds.client.Put(ds.ctx, newKey, &it)
  if err != nil {
    return "", err
  }

  return strconv.FormatInt(key.ID, 16), nil
}

func (ds *GoogleDataStore) Delete(key string) error {
  id, err := strconv.ParseInt(key, 16, 64)
  if err != nil {
    return err
  }

  ds_key := datastore.IDKey(TABLE, id, nil)

  // Load the key first to make sure it exists.
  var it Item
  if err := ds.client.Get(ds.ctx, ds_key, &it); err != nil {
    return err
  }

  return ds.client.Delete(ds.ctx, ds_key)
}

func (ds *GoogleDataStore) Update(key string, it Item) error {
  id, err := strconv.ParseInt(key, 16, 64)
  if err != nil {
    return err
  }
  ds_key := datastore.IDKey(TABLE, id, nil)

  // Load the key first to make sure it exists.
  var existing Item
  if err := ds.client.Get(ds.ctx, ds_key, &existing); err != nil {
    return err
  }

  _, err = ds.client.Put(ds.ctx, ds_key, &it)
  return err
}
