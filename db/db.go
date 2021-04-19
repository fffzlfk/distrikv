package db

import (
	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

// Database is an open bolt database
type Database struct {
	db *bolt.DB
}

// constructor
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	closeFunc = boltDb.Close

	db = &Database{boltDb}
	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, err
	}
	return
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// SetKey sets the key to the requested value or returns an error
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// SetKey gets the value of the requested from a default database
func (d *Database) GetKey(key string) (res []byte, err error) {
	err = d.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(defaultBucket)
		res = b.Get([]byte(key))
		return nil
	})
	return
}
