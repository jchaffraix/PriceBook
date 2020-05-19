package datastore

import (
  "context"
  "errors"
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

func encodeKeyId(key *datastore.Key) string {
  // TODO: Should I use Key.GobEncode/GobDecode?
  return strconv.FormatInt(key.ID, 16)
}

func decodeKeyId(key string) (int64, error) {
  if key == "" {
    return 0, errors.New("Missing key for update")
  }

  return strconv.ParseInt(key, 16, 64)
}

func (ds *GoogleDataStore) Add(userID string, it Item) (string, error) {
  err := ValidateItemWithoutKey(it)
  if err != nil {
    return "", err
  }

  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  newKey := datastore.IncompleteKey(ITEM_TABLE, ancestorKey)
  key, err := ds.client.Put(ds.ctx, newKey, &it)
  if err != nil {
    return "", err
  }

  return encodeKeyId(key), nil
}

func (ds *GoogleDataStore) Delete(userID, key string) error {
  keyID, err := decodeKeyId(key)
  if err != nil {
    return err
  }

  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  ds_key := datastore.IDKey(ITEM_TABLE, keyID, ancestorKey)

  // Load the key first to make sure it exists.
  var it Item
  if err := ds.client.Get(ds.ctx, ds_key, &it); err != nil {
    return err
  }

  return ds.client.Delete(ds.ctx, ds_key)
}

func (ds *GoogleDataStore) Update(userID string, it Item) error {
  err := ValidateItemWithKey(it)
  if err != nil {
    return err
  }

  keyID, err := decodeKeyId(it.ID)
  if err != nil {
    return err
  }
  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  dsKey := datastore.IDKey(ITEM_TABLE, keyID, ancestorKey)

  // Load the key first to make sure it exists.
  var existing Item
  if err := ds.client.Get(ds.ctx, dsKey, &existing); err != nil {
    return err
  }

  _, err = ds.client.Put(ds.ctx, dsKey, &it)
  return err
}

func (ds *GoogleDataStore) Get(userID string) []Item {
  ancestorKey := datastore.NameKey(USER_TABLE, userID, nil)
  query := datastore.NewQuery(ITEM_TABLE).Ancestor(ancestorKey)
  res := make([]Item, 0)
  it := ds.client.Run(ds.ctx, query)
  for {
    var item Item
    key, err := it.Next(&item)
    if err == iterator.Done {
      break;
    }
    if err != nil {
      log.Fatalf("Error reading the data: %+v", err)
    }

    // Sanity check.
    // If we fail, we just silently skip the item.
    err = ValidateItemWithoutKey(item)
    if err != nil {
      log.Printf("Item in storage doesn't pass validation: %+v", item)
      continue
    }

    // Populate the key as it's not automatically done.
    item.ID = encodeKeyId(key)
    res = append(res, item)
  }
  return res
}
