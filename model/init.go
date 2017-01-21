package model

import (
	"fmt"
	"github.com/boltdb/bolt"
)

// DB database handler
type DB struct {
	*bolt.DB
}

// OpenDB opens existing or initializes new db
func OpenDB(file string) (*DB, error) {
	db, err := bolt.Open(file, 0644, nil)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Close closes database
func (db *DB) Close() {
	if err := db.DB.Close(); err != nil {
		panic(err)
	}
}

// BucketFunc callback function for the newly created bucket
type BucketFunc func(bucket *bolt.Bucket) error

// WithNewBucket will create a new bucket and pass it to a given callback function
func WithNewBucket(tx *bolt.Tx, name string, callback BucketFunc) error {
	bucket, err := tx.CreateBucket([]byte(name))
	if err != nil {
		return err
	}

	callback(bucket)
	return nil
}

// Get load entry from bucket
func (db *DB) Get(bucket string, id string) ([]byte, error) {
	var result []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		_bucket := tx.Bucket([]byte(bucket))
		if _bucket == nil {
			return fmt.Errorf("Bucket %s was not found!", bucket)
		}

		result = _bucket.Get([]byte(id))
		if len(result) == 0 {
			result = nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
