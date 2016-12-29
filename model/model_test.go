package model

import (
  "os"
  "io/ioutil"
)

type TestDB struct {
  *DB
}

func NewTestDB() *TestDB {
  db, err := OpenDB(tempfile())
  if err != nil {
    panic(err)
  }

  return &TestDB{db}
}

func (db *TestDB) CloseAndDestroy() {
  defer os.Remove(db.Path())
  db.Close()
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
