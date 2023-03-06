package main

import "strings"

type Db map[string]any

func (db Db) Get(key string) any {
	return db[key]
}

func (db Db) Put(key string, value any) {
	db[key] = value
}

func (db Db) Delete(key string) {
	delete(db, key)
}

func (db Db) List(prefix string) []any {
	//goland:noinspection GoPreferNilSlice
	items := []any{}
	for key, item := range db {
		if strings.HasPrefix(key, prefix) {
			items = append(items, item)
		}
	}
	return items
}

var db = Db{}
