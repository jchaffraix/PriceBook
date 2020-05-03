package datastore

import (
  "context"
  "log"
  "strconv"

  "google.golang.org/api/iterator"
  "cloud.google.com/go/datastore"
)

const PROJECT_ID string = "pricebook"
const ITEM_TABLE string = "Items"
const USER_TABLE string = "Users"

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

func (ds *GoogleDataStore) Add(userID string, it Item) (string, error) {
  // TODO: Share the validation code with the InMemoryDataStore.
  if it.Name == "" {
    return "", &InvalidItemError{"Missing 'name'"}
  }

  // TODO: Pass the user instead of nil to scope their request.
  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  newKey := datastore.IncompleteKey(ITEM_TABLE, ancestorKey)
  key, err := ds.client.Put(ds.ctx, newKey, &it)
  if err != nil {
    return "", err
  }

  // TODO: Should I use Key.GobEncode/GobDecode?
  return strconv.FormatInt(key.ID, 16), nil
}

func (ds *GoogleDataStore) Delete(userID, key string) error {
  id, err := strconv.ParseInt(key, 16, 64)
  if err != nil {
    return err
  }

  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  ds_key := datastore.IDKey(ITEM_TABLE, id, ancestorKey)

  // Load the key first to make sure it exists.
  var it Item
  if err := ds.client.Get(ds.ctx, ds_key, &it); err != nil {
    return err
  }

  return ds.client.Delete(ds.ctx, ds_key)
}

func (ds *GoogleDataStore) Update(userID, key string, it Item) error {
  id, err := strconv.ParseInt(key, 16, 64)
  if err != nil {
    return err
  }
  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  ds_key := datastore.IDKey(ITEM_TABLE, id, ancestorKey)

  // Load the key first to make sure it exists.
  var existing Item
  if err := ds.client.Get(ds.ctx, ds_key, &existing); err != nil {
    return err
  }

  _, err = ds.client.Put(ds.ctx, ds_key, &it)
  return err
}

func (ds *GoogleDataStore) Get(userID string) []Item {
  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  query := datastore.NewQuery(ITEM_TABLE).Ancestor(ancestorKey)
  var res []Item
  it := ds.client.Run(ds.ctx, query)
  for {
    var item Item
    key, err := it.Next(&item)
    log.Printf("Key = %+v", key)
    if err == iterator.Done {
      break;
    }
    if err != nil {
      log.Fatalf("Error reading the data: %+v", err)
    }
    res = append(res, item)
  }
  return res
}
