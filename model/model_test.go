package model

import (
  "os"
  "io/ioutil"
)

// MustOpenDB return a new, open DN at a temporary location
func OpenTestDB() *DB {
  db, err := OpenDB(tempfile())
  if err != nil {
    panic(err)
  }

  return db
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
