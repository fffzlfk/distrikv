package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/fffzlfk/distrikv/utils"
)

type client struct {
	addr string
}

func (c *client) Set(key, value string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := fmt.Sprintf("http://%s/set?key=%s&value=%s", c.addr, key, value)
	req.SetRequestURI(url)
	fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	var respObj utils.Resp
	json.Unmarshal(bodyBytes, &respObj)
	if respObj.Err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("error")
	}
}

func (c *client) Del(key string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := fmt.Sprintf("http://%s/delete?key=%s", c.addr, key)
	req.SetRequestURI(url)
	fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	var respObj utils.Resp
	json.Unmarshal(bodyBytes, &respObj)
	if respObj.Err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("error")
	}
}

func (c *client) Get(key string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	url := fmt.Sprintf("http://%s/get?key=%s", c.addr, key)
	req.SetRequestURI(url)
	fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	var respObj utils.Resp
	json.Unmarshal(bodyBytes, &respObj)
	if respObj.Value == "" {
		fmt.Println(nil)
	} else {
		fmt.Printf("shard %d: %q\n", respObj.CurShard, respObj.Value)
	}
}

func (c *client) readInput() {
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
