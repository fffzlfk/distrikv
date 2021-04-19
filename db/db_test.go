package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fffzlfk/distrikv/db"
)

func TestGetSet(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "test.db")
	if err != nil {
		t.Fatal("could not create temp file:", err)
	}
	name := f.Name()
	f.Close()
	defer os.Remove(name)

	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatal("could not create a new database:", db)
	}
	defer closeFunc()

	if err := db.SetKey("setkey-test", []byte("good")); err != nil {
		t.Fatal("could not setkey:", err)
	}

	value, err := db.GetKey("setkey-test")
	if err != nil {
		t.Fatal(`could not get "setkey-test":`, err)
	}
	if !bytes.Equal(value, []byte("good")) {
		t.Fatalf(`unexpected value for key "setkey-test", got: %q, want: %q`, value, "good")
	}
}
