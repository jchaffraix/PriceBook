package datastore

type Item struct {
  // This is the hexadecimal representation of the Key.ID.
  // It is not stored but is sometimes returned to our API.
  // It is mandatory for updating and deleting.
  ID string `json:"id" datastore:"-"`
  Name string
  Quantity float32 `datastore:",noindex"`
  Unit string `datastore:",noindex"`

  // TODO: Need to add store
  // TODO: Which format for the date.
  // TODO: Add the timeseries.
}

type IDataStore interface {
  // TODO: Consolidate the add/update APIs into a single upsert?
  Add(userID string, it Item) (string, error)
  Delete(userID, key string) error
  Update(userID string, it Item) error
  Get(userID string) []Item
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
