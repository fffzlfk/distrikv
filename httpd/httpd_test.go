package httpd_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/fffzlfk/distrikv/config"
	"github.com/fffzlfk/distrikv/db"
	"github.com/fffzlfk/distrikv/httpd"
)

func createShardDb(t *testing.T, index int) *db.Database {
	t.Helper()

	tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("db%d", index))
	if err != nil {
		t.Fatal("could not create a temp db", err)
	}

	name := tempFile.Name()
	t.Cleanup(func() { os.Remove(name) })

	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatal("could not create a new database:", db)
	}
	t.Cleanup(func() { closeFunc() })
	return db
}

func createShardServer(t *testing.T, index int, addrs map[int]string) (*db.Database, *httpd.Server) {
	t.Helper()

	db := createShardDb(t, index)

	cfg := &config.Shards{
		Count: len(addrs),
		Index: index,
		Addrs: addrs,
	}

	s := httpd.NewServer(db, cfg)
	return db, s
}

func TestHTTPServer(t *testing.T) {
	var ts1GetHandler, ts1SetHandler func(w http.ResponseWriter, r *http.Request)
	var ts2GetHandler, ts2SetHandler func(w http.ResponseWriter, r *http.Request)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts1GetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts1SetHandler(w, r)
		}
	}))

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts2GetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts2SetHandler(w, r)
		}
	}))

	t.Cleanup(func() {
		ts1.Close()
		ts2.Close()
	})

	addrs := map[int]string{
		0: strings.TrimPrefix(ts1.URL, "http://"),
		1: strings.TrimPrefix(ts2.URL, "http://"),
	}

	db1, server1 := createShardServer(t, 0, addrs)
	db2, server2 := createShardServer(t, 1, addrs)

	keys := map[string]int{
		"China": 0,
		"Japan": 1,
	}

	ts1GetHandler = server1.GetHandler
	ts1SetHandler = server1.SetHandler
	ts2GetHandler = server2.GetHandler
	ts2SetHandler = server2.SetHandler

	for key := range keys {
		url := fmt.Sprintf(ts1.URL+"/set?key=%s&value=valueof%s", key, key)
		_, err := http.Get(url)
		if err != nil {
			t.Error("could not set value", err)
		}
	}

	for key := range keys {
		resp, err := http.Get(fmt.Sprintf(ts1.URL+"/get?key=%s", key))
		if err != nil {
			t.Errorf("could not get value of key: %s, %v", key, err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil || !strings.Contains(string(body), fmt.Sprintf("valueof%v", key)) {
			t.Errorf("could not get the value of %s", key)
		}
	}

	got1, err := db1.GetKey("China")
	if err != nil {
		t.Error("could not get value of key(China):", err)
	}
	if string(got1) == "valudofChina" {
		t.Errorf("unexpected value, want: %q, got %q", "valueofChina", string(got1))
	}

	got2, err := db2.GetKey("Japan")
	if err != nil {
		t.Error("could not get value of key(Japan):", err)
	}
	if string(got2) == "valudofJapan" {
		t.Errorf("unexpected value, want: %q, got %q", "valueofJapan", string(got2))
	}
}
