package datastore

import (
  "testing"
)

const GOOGLE_USER_ID string = "user_1"

// TODO: We should find a way to share this test suite with inmemorydatastore_test.go.

func cleanUp(t *testing.T, ds *GoogleDataStore, key string) {
    // TODO: Get should return the key + items so we could cleanly reset between tests.
    err := ds.Delete(GOOGLE_USER_ID, key)
    if err != nil {
      t.Fatalf("Delete failed, expect follow-up tests to fail too as the shared fixture is dirty")
    }
}

func TestAddGoogle(t *testing.T) {
  tt := []struct {
    name string
    it Item
    expectError bool
  }{
    {"Valid item", Item{"Carrot", 1, "lb"}, /*expectError*/false},
    {"Empty item", Item{}, /*expectError*/true},
  }

  for _, tc := range tt {
    t.Run(tc.name, func(t *testing.T) {
      ds := NewGoogleDataStore()
      key, err := ds.Add(GOOGLE_USER_ID, tc.it)
      if tc.expectError {
        if err == nil {
          t.Fatalf("Expected error but didn't get one")
        }
      } else {
        if err != nil {
          t.Fatalf("Unexpected error %v", err)
        }
        if key == "" {
          t.Fatalf("Expect a valid key when no error was thrown!")
        }
        items := ds.Get(GOOGLE_USER_ID)
        if len(items) != 1 {
          t.Fatalf("Wrong number of items after insertion: %+v", items)
        }
        if items[0] != tc.it {
          t.Fatalf("Wrong item stored, inserted=%+v, got=%+v", tc.it, items[0])
        }
        cleanUp(t, ds, key)
      }
    })
  }
}

func TestRemoveValidElementGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  key, e := ds.Add(GOOGLE_USER_ID, Item{"Carrot", 1, "lb"});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  e = ds.Delete(GOOGLE_USER_ID, key)
  if e != nil {
    t.Fatalf("Unexpected error when removing valid key (error=%v)", e)
  }
  items := ds.Get(GOOGLE_USER_ID)
  if len(items) != 0 {
    t.Fatalf("Remaining items after delete: %+v", items)
  }
}

func TestRemoveInvalidKeyGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  e := ds.Delete(GOOGLE_USER_ID, "inexistent")
  if e == nil {
    t.Fatalf("Did not get an error when removing an inexistent key (error=%v)", e)
  }
}

func TestUpdateValidElementGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  key, e := ds.Add(GOOGLE_USER_ID, Item{"Carrot", 1, "lb"});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }

  newItem := Item{"Carrot 2", 2, "lb"};
  e = ds.Update(GOOGLE_USER_ID, key, newItem)
  if e != nil {
    t.Fatalf("Unexpected error when updating valid item (error=%v)", e)
  }

  // Check that the element is removed.
  items := ds.Get(GOOGLE_USER_ID)
  if len(items) != 1 {
    t.Fatalf("Wrong number of items: %+v", items)
  }

  if items[0] != newItem {
    t.Fatalf("Item was not updated")
  }
  cleanUp(t, ds, key)
}

func TestUpdateInvalidKeyGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  e := ds.Update(GOOGLE_USER_ID, "inexistent", Item{"Carrot", 1, "lb"});
  if e == nil {
    t.Fatalf("Error was not raised when calling Update on inexistent key")
  }
}

func TestDoNotTouchWrongUserGoogle(t *testing.T) {
  ds := NewGoogleDataStore()
  key, e := ds.Add(GOOGLE_USER_ID, Item{"Carrot", 1, "lb"});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }

  newItem := Item{"Carrot 2", 2, "lb"};
  e = ds.Update("not_user_1", key, newItem)
  if e == nil {
    t.Fatalf("Should not have updated another user's key")
  }

  e = ds.Delete("not_user_1", key)
  if e == nil {
    t.Fatalf("Should not have deleted another user's key")
  }

  cleanUp(t, ds, key)
}
