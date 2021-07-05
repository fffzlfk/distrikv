package db

import (
	"bytes"
	"errors"

	bolt "go.etcd.io/bbolt"

	"github.com/fffzlfk/distrikv/utils"
)

// Database is an open bolt database
type Database struct {
	db       *bolt.DB
	readOnly bool
}

// constructor
func NewDatabase(dbPath string, readOnly bool) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	closeFunc = boltDb.Close

	db = &Database{boltDb, readOnly}
	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, err
	}
	return
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(t *bolt.Tx) error {
		if _, err := t.CreateBucketIfNotExists(utils.DefaultBucket); err != nil {
			return err
		}

		if _, err := t.CreateBucketIfNotExists(utils.ReplicaBucket); err != nil {
			return err
		}

		if _, err := t.CreateBucketIfNotExists(utils.DeleteBucket); err != nil {
			return err
		}
		return nil
	})
}

// SetKey sets the key to the requested value or returns an error
func (d *Database) SetKey(key string, value []byte) error {
	if d.readOnly {
		return errors.New("read only mode")
	}
	return d.db.Update(func(t *bolt.Tx) error {
		if err := t.Bucket(utils.DefaultBucket).Put([]byte(key), value); err != nil {
			return err
		}
		return t.Bucket(utils.ReplicaBucket).Put([]byte(key), value)
	})
}

// DeleteKey deletes the key to the requested value or returns an error
func (d *Database) DeleteKey(key string) error {
	// return d.SetKey(key, nil)
	return d.db.Update(func(t *bolt.Tx) error {
		value, err := d.GetKey(key)
		if err != nil {
			return err
		}
		if err := t.Bucket(utils.DefaultBucket).Delete([]byte(key)); err != nil {
			return err
		}
		return t.Bucket(utils.DeleteBucket).Put([]byte(key), value)
	})
}

// DeleteKeyOnReplica delete the key to the requested value into
// default databas for replicas
func (d *Database) DeleteKeyOnReplica(key string) error {
	return d.db.Update(func(t *bolt.Tx) error {
		return t.Bucket(utils.DefaultBucket).Delete([]byte(key))
	})
}

// SetKeyOnReplica set the key to the requested value into default database
// and does not write to the replication queue
// this method is only for replicas
func (d *Database) SetKeyOnReplica(key string, value []byte) error {
	return d.db.Update(func(t *bolt.Tx) error {
		return t.Bucket(utils.DefaultBucket).Put([]byte(key), value)
	})
}

// SetKey gets the value of the requested from a default database
func (d *Database) GetKey(key string) (res []byte, err error) {
	err = d.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(utils.DefaultBucket)
		res = b.Get([]byte(key))
		return nil
	})
	return
}

func copyByteSlice(src []byte) []byte {
	if src == nil {
		return nil
	}
	dest := make([]byte, len(src))
	copy(dest, src)
	return dest
}

// GetNextForReplication returns the key and value for the keys that have
// changed and have not yet been applied to replicas
func (d *Database) GetNextForReplicationOrDelete(bucket []byte) (key, value []byte, err error) {
	err = d.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucket)
		k, v := b.Cursor().First()
		key = copyByteSlice(k)
		value = copyByteSlice(v)
		return nil
	})

	if err != nil {
		key, value = nil, nil
	}
	return
}

// DeleteReplicationKey deletes the key from the replication queue
// if the value matches the contents or the key is already absent
func (d *Database) DeleteReplicationOrDeletedKey(bucket, key, value []byte) error {
	return d.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucket)

		v := b.Get(key)
		if v == nil {
			return errors.New("key does not exist")
		}

		if !bytes.Equal(v, value) {
			return errors.New("value does not match")
		}
		return b.Delete(key)
	})
}

// DeleteExtraKeys delete the keys that do not belongs to this shard
func (d *Database) DeleteExtraKeys(isExtra func(string) bool) error {
	var keys []string
	err := d.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(utils.DefaultBucket)
		return b.ForEach(func(k, v []byte) error {
			ks := string(k)
			if isExtra(ks) {
				keys = append(keys, ks)
			}
			return nil
		})
	})

	if err != nil {
		return err
	}

	return d.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(utils.DefaultBucket)

		for _, k := range keys {
			if err := b.Delete([]byte(k)); err != nil {
				return err
			}
		}
		return nil
	})
}
