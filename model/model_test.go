package model

import (
  "os"
  "io/ioutil"
)

type TestDB struct {
  *DB
  File string
}

// MustOpenDB return a new, open DN at a temporary location
func OpenTestDB() *TestDB {
  file := tempfile()
  db, err := OpenDB(file)
  if err != nil {
    panic(err)
  }

  return &TestDB{DB: db, File: file}
}

// CloseTestDB closes the database and deletes the underlying file; panic on error
func (db *TestDB) CloseTestDB() {
  if err := db.Close(); err != nil {
    panic(err)
  }
}

func tempfile() string {
  f, err := ioutil.TempFile("", "b4v-")
  if err != nil {
    panic(err)
  }

  if err := f.Close(); err != nil {
    panic(err)
  }

  if err := os.Remove(f.Name()); err != nil {
    panic(err)
  }

  return f.Name()
}
