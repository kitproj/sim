package main

// Db is a key-value store
type Db interface {
	Get(key string) any
	Put(key string, value any)
	Delete(key string)
	List(prefix string) []any
}

var db Db = diskDb("db")
