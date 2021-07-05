package replica

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

const (
	Replication = iota
	Deleted
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

func ClientLoop(db *db.Database, masterAddrs string, action int) {
	c := client{db: db, masterAddrs: masterAddrs}
	go func() {
		for {
			has, err := c.loop(action)
			if err != nil {
				log.Println("could not loop:", err)
				time.Sleep(time.Second)
				continue
			}

			if !has {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}

func (c *client) loop(action int) (bool, error) {
	var url string
	if action == Replication {
		url = "http://%s/next-replication-key"
	} else if action == Deleted {
		url = "http://%s/next-deleted-key"
	}

	url = fmt.Sprintf(url, c.masterAddrs)
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

	if action == Replication {
		if err := c.db.SetKeyOnReplica(res.Key, []byte(res.Value)); err != nil {
			return false, err
		}
		if err := c.deleteFromQueue(res.Key, res.Value, action); err != nil {
			log.Printf("could not deleteFromReplicationqueue(%q, %q): %v\n", res.Key, res.Value, err)
		}
	} else if action == Deleted {
		if err := c.db.DeleteKeyOnReplica(res.Key); err != nil {
			return false, err
		}
		if err := c.deleteFromQueue(res.Key, res.Value, action); err != nil {
			log.Printf("could not deleteFromDeletedqueue(%q, %q): %v\n", res.Key, res.Value, err)
		}
	}

	return true, nil
}

func (c *client) deleteFromQueue(key, value string, action int) error {
	u := url.Values{}
	u.Set("key", key)
	u.Set("value", value)

	var actionUrl string
	if action == Replication {
		actionUrl = "delete-replication-key"
	} else if action == Deleted {
		actionUrl = "delete-deleted-key"
	}

	log.Printf("deleting key=%q, value=%q from %s queue on %q", key, value, actionUrl, c.masterAddrs)

	url := fmt.Sprintf("http://%s/%s?%s", c.masterAddrs, actionUrl, u.Encode())

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
