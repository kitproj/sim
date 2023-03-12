package db

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type diskDb string

func (db diskDb) Get(key string) any {
	name := filepath.Join(string(db), key)
	file, err := os.Open(name)
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

func (db diskDb) Put(key string, value any) bool {
	name := filepath.Join(string(db), key)
	_ = os.MkdirAll(filepath.Dir(name), 0755)
	_, err := os.Stat(name)
	exists := !os.IsNotExist(err)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(value)
	if err != nil {
		panic(err)
	}
	return exists
}

func (db diskDb) Delete(key string) bool {
	name := filepath.Join(string(db), key)
	if err := os.Remove(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return true
}

func (db diskDb) List(prefix string) []any {
	//goland:noinspection GoPreferNilSlice
	items := []any{}
	err := filepath.WalkDir(filepath.Join(string(db), prefix), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			items = append(items, db.Get(strings.TrimPrefix(path, string(db+"/"))))
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	return items
}
