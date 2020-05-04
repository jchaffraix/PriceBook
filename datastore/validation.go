package datastore

func validateItemFields(item Item) error {
  if item.Name == "" {
    return &InvalidItemError{"Missing 'name'"}
  }
  return nil
}

// Validate that an item is valid and doesn't have a key.
//
// For Add, we assume that the ID field is not set
// as it should not be an existing entry.
// For Get, we should not store the key.
func ValidateItemWithoutKey(item Item) error {
  if item.ID != "" {
    return &InvalidItemError{"Unexpected ID"}
  }

  return validateItemFields(item)
}

// Validate that an item is valid and has a key.
// This is useful for update.
func ValidateItemWithKey(item Item) error {
  if item.ID == "" {
    return &InvalidItemError{"Missing ID"}
  }

  return validateItemFields(item)
}
