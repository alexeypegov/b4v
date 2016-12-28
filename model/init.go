package model

import (
	"github.com/boltdb/bolt"
)

// DB Bolt handler
var DB *bolt.DB

func init() {
	var err error
	if DB, err = bolt.Open("b4v.db", 0644, nil); err != nil {
		panic(err)
	}
}
