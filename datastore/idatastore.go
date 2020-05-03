package datastore

type Item struct {
  Name string
  Quantity float32
  Unit string
  // TODO: Which format for the date.
  // TODO: Add the timeseries.
}

type IDataStore interface {
  Add(it Item) (string, error)
  Delete(key string) error
  Update(key string, it Item) error
}

// Error raised when the Item to be stored is invalid.
type InvalidItemError struct {
  validationError string
}

func (e *InvalidItemError) Error() string {
  // TODO: Improve error reporting.
  return e.validationError
}

var ds IDataStore

func Init(isInMemory bool) {
  if isInMemory {
    ds = NewInMemoryDataStore();
  } else {
    ds = NewGoogleDataStore();
  }
}

func Get() *IDataStore {
  return &ds;
}
