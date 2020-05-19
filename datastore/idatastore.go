package datastore


// PurchaseInfo consistute one point in
// a time series of purchases.
type PurchaseInfo struct {
  // TODO: We should probably use time.Time.
  // Unfortunately this is not parsed correctly from
  // our front-end ISO date.
  Date string `json:"date"`
  Store string `json:"store"`
  Price float32 `json:"price,string"`
  Currency string `json:"currency"`
}

type Item struct {
  // This is the hexadecimal representation of the Key.ID.
  // It is not stored but is sometimes returned to our API.
  // It is mandatory for updating and deleting.
  ID string `json:"id, string" datastore:"-"`
  Name string `json:"name" datastore:",noindex"`
  Brand string `json:"brand" datastore:",noindex"`
  Quantity int32 `json:"quantity,string" datastore:",noindex"`
  Unit string `json:"unit" datastore:",noindex"`

  Purchases []PurchaseInfo `json:"purchases"`
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
