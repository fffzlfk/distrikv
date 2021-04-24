package replication

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/fffzlfk/distrikv/db"
)

type NextKeyValue struct {
	Key   string
	Value string
	Err   error
}

type client struct {
	db          *db.Database
	masterAddrs string
}

func ClientLoop(db *db.Database, masterAddrs string) {
	c := client{db: db, masterAddrs: masterAddrs}
	for {
		has, err := c.loop()
		if err != nil {
			log.Println("could not loop:", err)
			time.Sleep(time.Second)
			continue
		}

		if !has {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (c *client) loop() (bool, error) {
	url := fmt.Sprintf("http://%s/next-replication-key", c.masterAddrs)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var res NextKeyValue
	json.NewDecoder(resp.Body).Decode(&res)
	if res.Err != nil {
		return false, err
	}

	if res.Key == "" {
		return false, nil
	}

	if err := c.db.SetKeyOnReplica(res.Key, []byte(res.Value)); err != nil {
		return false, err
	}
	if err := c.deleteFromReplicationQueue(res.Key, res.Value); err != nil {
		log.Printf("could not deleteFromReplicationqueue(%q, %q): %v\n", res.Key, res.Value, err)
	}
	return true, nil
}

func (c *client) deleteFromReplicationQueue(key, value string) error {
	u := url.Values{}
	u.Set("key", key)
	u.Set("value", value)

	log.Printf("deleting key=%q, value=%q from replication queue on %q", key, value, c.masterAddrs)

	url := fmt.Sprintf("http://%s/delete-replication-key?%s", c.masterAddrs, u.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !bytes.Equal(rb, []byte("ok")) {
		return errors.New(string(rb))
	}
	return nil
}
