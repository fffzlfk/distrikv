package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fffzlfk/distrikv/db"
	"github.com/fffzlfk/distrikv/utils"
)

func setKey(t *testing.T, d *db.Database, key, value string) {
	t.Helper()
	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("could not Setkey(%q, %q): %v", key, value, err)
	}
}

func getKey(t *testing.T, d *db.Database, key string) string {
	t.Helper()
	res, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("could not Getkey(%q): %v", key, err)
	}
	return string(res)
}

func delKey(t *testing.T, d *db.Database, key string) {
	t.Helper()
	err := d.DeleteKey(key)
	if err != nil {
		t.Fatalf("could not Getkey(%q): %v", key, err)
	}
}

func createTempDb(t *testing.T, readOnly bool) *db.Database {
	t.Helper()
	f, err := ioutil.TempFile(os.TempDir(), "test.db")
	if err != nil {
		t.Fatal("could not create temp file:", err)
	}
	name := f.Name()
	t.Cleanup(func() {
		f.Close()
		os.Remove(name)
	})

	db, closeFunc, err := db.NewDatabase(name, readOnly)
	if err != nil {
		t.Fatal("could not create a new database:", db)
	}
	t.Cleanup(func() {
		err := closeFunc()
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}

func TestGetSet(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "setkey-test", "good")
	setKey(t, db, "setkey-extratest", "good")

	if value := getKey(t, db, "setkey-test"); value != "good" {
		t.Fatalf(`unexpected value for key "setkey-test", got: %q, want: %q`, value, "good")
	}

	if err := db.DeleteExtraKeys(func(s string) bool { return s == "setkey-extratest" }); err != nil {
		t.Fatalf(`coult not DeleteExtraKeys("setkey-test"): %v`, err)
	}

	if value := getKey(t, db, "setkey-test"); value != "good" {
		t.Fatalf(`unexpected value for key "setkey-test", got: %q, want: %q`, value, "good")
	}

	if value := getKey(t, db, "setkey-extratest"); value != "" {
		t.Fatalf(`unexpected value for key "setkey-extratest", got: %q, want: %q`, value, "")
	}
}

func TestDeleteReplicationKey(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "setkey-test", "good")

	k, v, err := db.GetNextForReplicationOrDelete(utils.ReplicaBucket)
	if err != nil {
		t.Fatal("could not GetNextForReplication:", err)
	}

	if !bytes.Equal(k, []byte("setkey-test")) || !bytes.Equal(v, []byte("good")) {
		t.Fatalf(`GetNextForReplication(): got %q, %q; want %q %q`, k, v, "setkey-test", "good")
	}

	if err := db.DeleteReplicationOrDeletedKey(utils.ReplicaBucket, k, v); err != nil {
		t.Fatal("could not DeleteReplicationKey:", err)
	}

	k, v, err = db.GetNextForReplicationOrDelete(utils.ReplicaBucket)
	if err != nil {
		t.Fatal("could not GetNextForReplication:", err)
	}

	if k != nil || v != nil {
		t.Fatalf(`GetNextForReplication(): got %q, %q; want nil nil`, k, v)
	}
}

func TestDeleteDeletedKey(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "setkey-test", "good")

	delKey(t, db, "setkey-test")

	k, v, err := db.GetNextForReplicationOrDelete(utils.DeleteBucket)
	if err != nil {
		t.Fatal("could not GetNextForReplication:", err)
	}

	if !bytes.Equal(k, []byte("setkey-test")) || !bytes.Equal(v, []byte("good")) {
		t.Fatalf(`GetNextForDeleted(): got %q, %q; want %q %q`, k, v, "setkey-test", "good")
	}

	if err := db.DeleteReplicationOrDeletedKey(utils.DeleteBucket, k, v); err != nil {
		t.Fatal("could not DeleteDeletedKey:", err)
	}

	k, v, err = db.GetNextForReplicationOrDelete(utils.DeleteBucket)
	if err != nil {
		t.Fatal("could not GetNextForDeleted:", err)
	}

	if k != nil || v != nil {
		t.Fatalf(`GetNextForDeleted(): got %q, %q; want nil nil`, k, v)
	}
}

func TestSetReadOnly(t *testing.T) {
	tmpDb := createTempDb(t, true)

	if err := tmpDb.SetKey("setkey-test", []byte("good")); err == nil {
		t.Fatalf("Setkey(%q, %q), got: nil err, want: not nil err", "setkry-test", "good")
	}
}
