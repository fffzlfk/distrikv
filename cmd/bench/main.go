package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	iterations  = flag.Int("iterations", 1000, "set the number of iterations for writing")
	concurrency = flag.Int("concurrency", 1, "set the number of goroutines")
)

func benchmark(name string, fn func()) {
	var min, max time.Duration
	min = time.Hour
	start := time.Now()
	for i := 0; i < *iterations; i++ {
		itStart := time.Now()
		fn()
		itTime := time.Since(itStart)
		if itTime < min {
			min = itTime
		}
		if itTime > max {
			max = itTime
		}
	}
	avg := time.Since(start) / (time.Duration(*iterations))
	qps := float64(*iterations) / (float64(time.Since(start)) / float64((time.Second)))
	fmt.Printf("func %s took time: avg=%v, QPS=%.1f, min=%v, max=%v\n", name, avg, qps, min, max)
}

func writeRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(10000))
	value := fmt.Sprintf("value-%d", rand.Intn(10000))

	values := url.Values{}

	values.Set("key", key)
	values.Set("value", value)

	resp, err := http.Get("http://localhost:8080" + "/set?" + values.Encode())
	if err != nil {
		log.Fatal("could not set:", err)
	}
	defer resp.Body.Close()
}

func init() {
	flag.Parse()
}

func main() {
	rand.Seed(int64(time.Now().UnixNano()))
	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			benchmark("write", writeRand)
			wg.Done()
		}()
	}
	wg.Wait()
}
