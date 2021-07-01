package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
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
	fmt.Println(string(bodyBytes))
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
	fmt.Println(string(bodyBytes))
}

func (c *client) readInput() {
	for {
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
		}
	}
}

func main() {
	// c := &client{"localhost:8021"}
	// c.Get("2")

	// ch := make(chan string)

}
