package datastore

import (
  "testing"
)

const GOOGLE_USER_ID string = "user_1"

// TODO: We should find a way to share this test suite with inmemorydatastore_test.go.

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
      key, e := ds.Add(GOOGLE_USER_ID, tc.it)
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
  // TODO: Add a Query to the API to validate that it is gone.
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
  // TODO: Add a Query to the API to validate that it is gone.
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
}
