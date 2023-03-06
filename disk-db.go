package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type diskDb string

func (db diskDb) Get(key string) any {
	file, err := os.Open(filepath.Join(string(db), key))
	if err != nil {
		return nil
	}
	defer file.Close()
	var value any
	err = json.NewDecoder(file).Decode(&value)
	if err != nil {
		panic(err)
	}
	return value
}

func (db diskDb) Put(key string, value any) {
	name := filepath.Join(string(db), key)
	_ = os.MkdirAll(filepath.Dir(name), 0755)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(value)
	if err != nil {
		panic(err)
	}
}

func (db diskDb) Delete(key string) {
	_ = os.Remove(filepath.Join(string(db), key))
}

func (db diskDb) List(prefix string) []any {
	//goland:noinspection GoPreferNilSlice
	items := []any{}
	err := filepath.WalkDir(filepath.Join(string(db), prefix), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			items = append(items, db.Get(strings.TrimPrefix(path, "db/")))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return items
}
