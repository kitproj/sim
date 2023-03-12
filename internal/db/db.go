package db

import (
	"os"
	"path/filepath"
)

// Db is a key-value store
type Db interface {
	Get(key string) any
	Put(key string, value any) bool
	Delete(key string) bool
	List(prefix string) []any
}

var Instance Db = diskDb(filepath.Join(os.Getenv("HOME"), ".kitproj", "sim", "db"))
