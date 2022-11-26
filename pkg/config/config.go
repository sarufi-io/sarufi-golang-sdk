package config

type Reader interface {
	// String returns the value of the environment variable
	// with the given key. If the key is not found, it returns
	// the default value
	String(key string, defaultValue string) string

	// Int returns the value of the environment variable
	// with the given key. If the key is not found, it returns
	// the default value
	Int(key string, defaultValue int64) int64

	// Bool returns the value of the environment variable
	// with the given key. If the key is not found, it returns
	// the default value
	Bool(key string, defaultValue bool) bool

	// Float returns the value of the environment variable
	// with the given key. If the key is not found, it returns
	// the default value
	Float(key string, defaultValue float64) float64

	// Add insert a new key value pair to the environment that
	// is not present in a list of currently loaded environment
	// variables
	Add(key string, value string) error

	// Update updates the value of an existing key in the environment
	// if the key is not present in the list of currently loaded
	// it returns an error
	Update(key string, value string) error

	// Delete deletes a key from the environment if the key is not
	// present in the list of currently loaded environment variables
	// it returns an error
	Delete(key string) error
}
