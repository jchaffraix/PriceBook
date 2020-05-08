package datastore

import (
  "time"
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

