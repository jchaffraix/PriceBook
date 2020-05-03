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
