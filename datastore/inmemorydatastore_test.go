package datastore

import (
  "testing"
)

const INMEMORY_USER_ID string = "user-123"

func TestAddInMemory(t *testing.T) {
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
      ds := NewInMemoryDataStore()
      key, e := ds.Add(INMEMORY_USER_ID, tc.it)
      if tc.expectError {
        if e == nil {
          t.Fatalf("Expected error but didn't get one")
        }
      } else {
        if e != nil {
          t.Fatalf("Unexpected error %v", e)
        }
        if key == "" {
          t.Fatalf("Expect a valid key when no error was thrown!")
        }
        items := ds.Get(INMEMORY_USER_ID)
        if len(items) != 1 {
          t.Fatalf("Wrong number of items after insertion: %+v", items)
        }
        if items[0] != tc.it {
          t.Fatalf("Wrong item stored, inserted=%+v, got=%+v", tc.it, items[0])
        }
      }
    })
  }
}

func TestRemoveValidElementInMemory(t *testing.T) {
  ds := NewInMemoryDataStore()
  key, e := ds.Add(INMEMORY_USER_ID, Item{"Carrot", 1, "lb"});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }
  e = ds.Delete(INMEMORY_USER_ID, key)
  if e != nil {
    t.Fatalf("Unexpected error when removing valid key (error=%v", e)
  }
  items := ds.Get(INMEMORY_USER_ID)
  if len(items) != 0 {
    t.Fatalf("Remaining items after delete: %+v", items)
  }
}

func TestRemoveInvalidKeyInMemory(t *testing.T) {
  ds := NewInMemoryDataStore()
  e := ds.Delete(INMEMORY_USER_ID, "inexistent")
  if e == nil {
    t.Fatalf("Did not get an error when removing an inexistent key")
  }
}

func TestUpdateValidElementInMemory(t *testing.T) {
  ds := NewInMemoryDataStore()
  key, e := ds.Add(INMEMORY_USER_ID, Item{"Carrot", 1, "lb"});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item (error=%v)", e)
  }

  newItem := Item{"Carrot 2", 2, "lb"};
  e = ds.Update(INMEMORY_USER_ID, key, newItem)
  if e != nil {
    t.Fatalf("Unexpected error when updating valid item (error=%v)", e)
  }

  // Check that the element is removed.
  items := ds.Get(INMEMORY_USER_ID)
  if len(items) != 1 {
    t.Fatalf("Wrong number of items: %+v", items)
  }

  if items[0] != newItem {
    t.Fatalf("Item was not updated")
  }
}

func TestUpdateInvalidKeyInMemory(t *testing.T) {
  ds := NewInMemoryDataStore()
  e := ds.Update(INMEMORY_USER_ID, "inexistent", Item{"Carrot", 1, "lb"});
  if e == nil {
    t.Fatalf("Error was not raised when calling Update on inexistent key")
  }
}

func TestDoNotTouchWrongUserInMemory(t *testing.T) {
  ds := NewInMemoryDataStore()
  key, e := ds.Add(INMEMORY_USER_ID, Item{"Carrot", 1, "lb"});
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
}
