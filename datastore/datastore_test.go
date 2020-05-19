package datastore

import (
  "time"
  "testing"
)

const USER_ID string = "user_1"

func createPurchaseInfo(t time.Time, store string, price float32, currency string) []PurchaseInfo {
  return []PurchaseInfo{
    PurchaseInfo{t, store, price, currency},
  }
}

func purchasesAreEqual(p1, p2 PurchaseInfo) bool {
  // We have to get the actual timestamp from the purchase.
  // This is because time.Time contains a monotonic TS by
  // default that confuses a strict equality.
  if p1.Time.Unix() != p2.Time.Unix() {
    return false;
  }
  if p1.Store != p2.Store {
    return false;
  }
  if p1.Currency != p2.Currency {
    return false
  }
  if p1.Price != p2.Price {
    return false;
  }
  return true;
}

func itemsAreEqual(i1, i2 Item) bool {
  if i1.ID != i2.ID {
    return false;
  }

  if i1.Name != i2.Name {
    return false;
  }

  if i1.Quantity != i2.Quantity {
    return false;
  }

  if i1.Unit != i2.Unit {
    return false;
  }

  if len(i1.Purchases) != len(i2.Purchases) {
    return false;
  }

  for i, purchase := range i1.Purchases {
    if !purchasesAreEqual(purchase, i2.Purchases[i]) {
      return false;
    }
  }

  return true;
}

// Actual tests.
// Hoisted here to be shared between the 2 implementations.
type DataStoreFactory func() IDataStore

func dataStoreCleanUp(t *testing.T, ds IDataStore, key string) {
    err := ds.Delete(USER_ID, key)
    if err != nil {
      t.Fatalf("Delete failed, expect follow-up tests to fail too as the shared fixture is dirty")
    }
}

func testAdd(t *testing.T, factory DataStoreFactory) {
  tt := []struct {
    name string
    it Item
    expectError bool
  }{
    {"Valid item without brand", Item{/*ID=*/"", "Yogurt", "Foo", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")}, /*expectError*/false},
    {"Valid item with brand", Item{/*ID=*/"", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")}, /*expectError*/false},
    {"Valid item with brand and decimal", Item{/*ID=*/"", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 2.99, "EUR")}, /*expectError*/false},
    {"Item without Name", Item{/*ID=*/"", /*Name=*/"", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")}, /*expectError*/true},
    {"Item without PurchaseInfo", Item{/*ID=*/"", "Carrot", "", 1, "lb", []PurchaseInfo{}}, /*expectError*/true},
    {"Item with a key", Item{/*ID=*/"1234", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")}, /*expectError*/true},
  }

  for _, tc := range tt {
    t.Run(tc.name, func(t *testing.T) {
      ds := factory()
      key, err := ds.Add(USER_ID, tc.it)
      if tc.expectError {
        if err == nil {
          t.Fatalf("Expected error but didn't get one")
        }
      } else {
        defer dataStoreCleanUp(t, ds, key)
        if err != nil {
          t.Fatalf("Unexpected error %v", err)
        }
        if key == "" {
          t.Fatalf("Expect a valid key when no error was thrown!")
        }
        items := ds.Get(USER_ID)
        if len(items) != 1 {
          t.Fatalf("Wrong number of items after insertion: %+v", items)
        }

        // Add the key to the item in Get.
        expectedItem := tc.it
        expectedItem.ID = key
        if !itemsAreEqual(items[0], expectedItem) {
          t.Fatalf("Wrong item stored, expected=%+v, got=%+v", expectedItem, items[0])
        }
      }
    })
  }
}

func testRemoveValidElement(t *testing.T, factory DataStoreFactory) {
  ds := factory()
  key, e := ds.Add(USER_ID, Item{"", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")});
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

func testRemoveInvalidKey(t *testing.T, factory DataStoreFactory) {
  ds := factory()
  e := ds.Delete(USER_ID, "inexistent")
  if e == nil {
    t.Fatalf("Did not get an error when removing an inexistent key (error=%v)", e)
  }
}

func testUpdateValidElement(t *testing.T, factory DataStoreFactory) {
  ds := factory()
  key, e := ds.Add(USER_ID, Item{"", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  defer dataStoreCleanUp(t, ds, key)

  newItem := Item{key, "Carrot 2", "", 2, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")};
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

func testUpdateInvalidKey(t *testing.T, factory DataStoreFactory) {
  ds := factory()
  e := ds.Update(USER_ID, Item{"inexistent", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")});
  if e == nil {
    t.Fatalf("Error was not raised when calling Update on inexistent key")
  }
}

func testUpdateInvalidElement(t *testing.T, factory DataStoreFactory) {
  ds := factory()
  originalItem := Item{"", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")}
  key, e := ds.Add(USER_ID, originalItem);
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  defer dataStoreCleanUp(t, ds, key)

  // Missing name.
  newItem := Item{key, "", "", 2,"lb", createPurchaseInfo(time.Now(), "Location", 42, "$")};
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


func testDoNotTouchWrongUser(t *testing.T, factory DataStoreFactory) {
  ds := factory()
  key, e := ds.Add(USER_ID, Item{"", "Carrot", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  defer dataStoreCleanUp(t, ds, key)

  newItem := Item{key, "Carrot 2", "", 2, "lb", createPurchaseInfo(time.Now(), "Location", 42, "$")};
  e = ds.Update("not_user_1", newItem)
  if e == nil {
    t.Fatalf("Should not have updated another user's key")
  }

  e = ds.Delete("not_user_1", key)
  if e == nil {
    t.Fatalf("Should not have deleted another user's key")
  }
}
