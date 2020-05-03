package datastore

import (
  "testing"
)

func TestAdd(t *testing.T) {
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
      key, e := ds.Add(tc.it)
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

func TestRemoveValidElement(t *testing.T) {
  ds := NewInMemoryDataStore()
  key, e := ds.Add(Item{"Carrot", 1, "lb"});
  if e != nil {
    t.Fatalf("Unexpected error when inserting valid item")
  }
  e = ds.Delete(key)
  if e != nil {
    t.Fatalf("Unexpected error when removing valid key")
  }
  _, found := ds.m[key]
  if found {
    t.Fatalf("Key was not removed despite no error returned")
  }
}

func TestRemoveInvalidKey(t *testing.T) {
  ds := NewInMemoryDataStore()
  e := ds.Delete("inexistent")
  if e == nil {
    t.Fatalf("Did not get an error when removing an inexistent key")
  }
}
