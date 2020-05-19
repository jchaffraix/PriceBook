package datastore

import (
  "testing"
)


// TODO: Not sure why I need this explicit wrapper.
func testGoogleDataStore() IDataStore {
  return NewGoogleDataStore();
}

func TestAddGoogle(t *testing.T) {
  testAdd(t, testGoogleDataStore)
}

func TestRemoveValidElementGoogle(t *testing.T) {
  testRemoveValidElement(t, testGoogleDataStore)
}

func TestRemoveInvalidKeyGoogle(t *testing.T) {
  testRemoveInvalidKey(t, testGoogleDataStore)
}

func TestUpdateValidElementGoogle(t *testing.T) {
  testUpdateValidElement(t, testGoogleDataStore)
}

func TestUpdateInvalidKeyGoogle(t *testing.T) {
  testUpdateInvalidKey(t, testGoogleDataStore)
}

func TestUpdateInvalidElementGoogle(t *testing.T) {
  testUpdateInvalidElement(t, testGoogleDataStore)
}

func TestDoNotTouchWrongUserGoogle(t *testing.T) {
  testDoNotTouchWrongUser(t, testGoogleDataStore)
}
