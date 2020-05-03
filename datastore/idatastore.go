package datastore

type Item struct {
  Name string
  Quantity float32
  Unit string
  // TODO: Which format for the date.
  // TODO: Add the timeseries.
}

type IDataStore interface {
  Add(userID string, it Item) (string, error)
  Delete(userID string, key string) error
  Update(userID string, key string, it Item) error
}

// Error raised when the Item to be stored is invalid.
type InvalidItemError struct {
  validationError string
}

func (e *InvalidItemError) Error() string {
  // TODO: Improve error reporting.
  return e.validationError
}

// Error raised when an Item is not found in the datastore.
type NotFoundError struct {
  key string
}

func (e *NotFoundError) Error() string {
  return e.key + " was not found"
}

var ds IDataStore

func Init(isInMemory bool) {
  if isInMemory {
    ds = NewInMemoryDataStore();
  } else {
    ds = NewGoogleDataStore();
  }
}

func Get() IDataStore {
  return ds;
}
