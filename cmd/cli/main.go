package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/fffzlfk/distrikv/utils"
)

type client struct {
	addr string
}

func (c *client) do(url string) (*utils.Resp, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(url)
	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}

	bodyBytes := resp.Body()
	var respObj utils.Resp
	err = json.Unmarshal(bodyBytes, &respObj)
	if err != nil {
		return nil, err
	}
	if respObj.Err != nil {
		return nil, respObj.Err
	}
	return &respObj, nil
}

func (c *client) Ping() error {
	url := fmt.Sprintf("http://%s/ping", c.addr)
	_, err := c.do(url)
	return err
}

func (c *client) Set(key, value string) {
	url := fmt.Sprintf("http://%s/set?key=%s&value=%s", c.addr, key, value)

	if _, err := c.do(url); err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("error")
	}
}

func (c *client) Del(key string) {
	url := fmt.Sprintf("http://%s/delete?key=%s", c.addr, key)
	if _, err := c.do(url); err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("error")
	}
}

func (c *client) Get(key string) {
	url := fmt.Sprintf("http://%s/get?key=%s", c.addr, key)
	resp, err := c.do(url)
	if err != nil {
		fmt.Println("error")
		return
	}
	if resp.Value == "" {
		fmt.Println(nil)
	} else {
		fmt.Printf("shard %d: %q\n", resp.CurShard, resp.Value)
	}
}

func (c *client) readInput() {
	if err := c.Ping(); err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Printf("%s> ", c.addr)
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			continue
		}
		msg = strings.ToLower(msg)
		msg = strings.Trim(msg, "\n")
		args := strings.Split(msg, " ")
		cmd := strings.Trim(args[0], " ")

		switch cmd {
		case "set":
			if len(args) < 3 {
				fmt.Println("expected 2 args")
				continue
			}
			c.Set(args[1], args[2])
		case "get":
			if len(args) < 2 {
				fmt.Println("expected 2 args")
				continue
			}
			c.Get(args[1])
		case "del":
			if len(args) < 2 {
				fmt.Println("expected 2 args")
				continue
			}
			c.Del(args[1])
		default:
			fmt.Println("syntax error")
		}
	}
}

func main() {
	c := &client{"localhost:8011"}
	c.readInput()
}
