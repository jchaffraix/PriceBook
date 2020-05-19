package datastore

import (
  "time"
  "testing"
)


// TODO: We should find a way to share this test suite with inmemorydatastore_test.go.

// TODO: Not sure why I need this explicit wrapper.
func testGoogleDataStore() IDataStore {
  return NewGoogleDataStore();
}

func TestAddGoogle(t *testing.T) {
  testAdd(t, testGoogleDataStore)
}

func TestRemoveValidElementGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  key, e := ds.Add(USER_ID, Item{"", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  e = ds.Delete(USER_ID, key)
  if e != nil {
    t.Fatalf("Unexpected error when removing valid key (error=%v)", e)
  }
  items := ds.Get(USER_ID)
  if len(items) != 0 {
    t.Fatalf("Remaining items after delete: %+v", items)
  }
}

func TestRemoveInvalidKeyGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  e := ds.Delete(USER_ID, "inexistent")
  if e == nil {
    t.Fatalf("Did not get an error when removing an inexistent key (error=%v)", e)
  }
}

func TestUpdateValidElementGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  key, e := ds.Add(USER_ID, Item{"", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  defer dataStoreCleanUp(t, ds, key)

  newItem := Item{key, "Carrot 2", 2, "lb", createPurchaseInfo(time.Now(), "Location", 42)};
  e = ds.Update(USER_ID, newItem)
  if e != nil {
    t.Fatalf("Unexpected error when updating valid item (error=%v)", e)
  }

  // Check that the element is removed.
  items := ds.Get(USER_ID)
  if len(items) != 1 {
    t.Fatalf("Wrong number of items: %+v", items)
  }

  if !itemsAreEqual(items[0], newItem) {
    t.Fatalf("Item was not updated")
  }
}

func TestUpdateInvalidKeyGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  e := ds.Update(USER_ID, Item{"inexistent", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)});
  if e == nil {
    t.Fatalf("Error was not raised when calling Update on inexistent key")
  }
}

func TestUpdateInvalidElementGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  originalItem := Item{"", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)}
  key, e := ds.Add(USER_ID, originalItem);
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  defer dataStoreCleanUp(t, ds, key)

  // Missing name.
  newItem := Item{key, "", 2, "lb", createPurchaseInfo(time.Now(), "Location", 42)};
  e = ds.Update(USER_ID, newItem)
  if e == nil {
    t.Fatalf("Expected error when updating with an invalid item but didn't get it!")
  }

  items := ds.Get(USER_ID)
  if len(items) != 1 {
    t.Fatalf("Unexpected item count: %+v", items)
  }
  originalItem.ID = key
  if !itemsAreEqual(items[0], originalItem) {
    t.Fatalf("Unexpected difference in storage. Expected %+v, got: %+v", originalItem, items[0])
  }
}


func TestDoNotTouchWrongUserGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  key, e := ds.Add(USER_ID, Item{"", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  defer dataStoreCleanUp(t, ds, key)

  newItem := Item{key, "Carrot 2", 2, "lb", createPurchaseInfo(time.Now(), "Location", 42)};
  e = ds.Update("not_user_1", newItem)
  if e == nil {
    t.Fatalf("Should not have updated another user's key")
  }

  e = ds.Delete("not_user_1", key)
  if e == nil {
    t.Fatalf("Should not have deleted another user's key")
  }
}
