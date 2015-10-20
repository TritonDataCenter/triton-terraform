package helpers

import (
	"os"
)

// EnvOrElse gets the environment variable specified or the default
func EnvOrElse(key, def string) string {
	if maybe := os.Getenv(key); maybe != "" {
		return maybe
	}

	return def
}

// SetUnset sets a value in the environment and returns a func that can be used
// to defer unsetting it
func SetUnset(key, value string) func() {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
	return func() {
		os.Unsetenv(key)
	}
}
