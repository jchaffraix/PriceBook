package datastore

import (
  "testing"
)


// TODO: Not sure why I need this explicit wrapper.
func testInMemoryDataStore() IDataStore {
  return NewInMemoryDataStore();
}

func TestAddInMemory(t *testing.T) {
  testAdd(t, testInMemoryDataStore)
}

func TestRemoveValidElementInMemory(t *testing.T) {
  testRemoveValidElement(t, testInMemoryDataStore)
}

func TestRemoveInvalidKeyInMemory(t *testing.T) {
  testRemoveInvalidKey(t, testInMemoryDataStore)
}

func TestUpdateValidElementInMemory(t *testing.T) {
  testUpdateValidElement(t, testInMemoryDataStore)
}

func TestUpdateInvalidElementInMemory(t *testing.T) {
  testUpdateInvalidElement(t, testInMemoryDataStore)
}

func TestUpdateInvalidKeyInMemory(t *testing.T) {
  testUpdateInvalidKey(t, testInMemoryDataStore)
}

func TestDoNotTouchWrongUserInMemory(t *testing.T) {
  testDoNotTouchWrongUser(t, testInMemoryDataStore)
}
