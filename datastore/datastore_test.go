package datastore

import (
  "time"
  "testing"
)

const USER_ID string = "user_1"

func createPurchaseInfo(t time.Time, store string, price float32) []PurchaseInfo {
  return []PurchaseInfo{
    PurchaseInfo{t, store, price},
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
    {"Valid item", Item{/*ID=*/"", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)}, /*expectError*/false},
    {"Item without Name", Item{/*ID=*/"", "", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)}, /*expectError*/true},
    {"Item without PurchaseInfo", Item{/*ID=*/"", "Carrot", 1, "lb", []PurchaseInfo{}}, /*expectError*/true},
    {"Item with a key", Item{/*ID=*/"1234", "Carrot", 1, "lb", createPurchaseInfo(time.Now(), "Location", 42)}, /*expectError*/true},
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
