package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	connectErrors  int32
	not200Errors   int32
	notFoundErrors int32
	success        int32
	reqCount       int32
	requestTime    int64
	startTime      time.Time
)

var (
	urlChan chan string

	reqWg    sync.WaitGroup
	workerWg sync.WaitGroup

	unrestricted bool
	filename     string
	numWorkers   int
	verbose      bool
	timeout      time.Duration
)

func stats() {
	if unrestricted {
		fmt.Println("Unrestricted")
	} else {
		fmt.Println("Concurrent Requests =", numWorkers)
	}

	fmt.Println("Number of Requests =", reqCount)
	fmt.Println("Connection Errors =", connectErrors)
	fmt.Println("Not 200 =", not200Errors)
	fmt.Println("Not Found (404) =", notFoundErrors)
	fmt.Println("Success =", success)
	fmt.Println("Avg Request Time (ms) =", float64(requestTime)/float64(reqCount)/1000000)
	fmt.Println("Time Taken (sec) =", time.Since(startTime))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.BoolVar(&unrestricted, "unrestricted", false, "make all requests with no restrictions on concurrency")
	flag.IntVar(&numWorkers, "concurrent", 0, "number of concurrent requests to make if not unrestricted")
	flag.StringVar(&filename, "file", "", "file of urls to process")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "http request timeout")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.Parse()

	// if running unrestricted, we only need a single worker
	if unrestricted {
		numWorkers = 1
	} else if numWorkers == 0 {
		fmt.Println("Must specify wither --unrestricted or number of concurrent requests (--concurrent).")
		return
	}

	startTime = time.Now().UTC()

	urlChan = make(chan string, 1000)

	defer stats()
	go func() {
		ticker := time.Tick(10 * time.Second)
		first := true
		for _ = range ticker {
			stats()
			if first {
				fmt.Println("==================================")
				first = false
			}
		}
	}()

	for x := 0; x < numWorkers; x++ {
		workerWg.Add(1)
		go processIds()
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var idx int
	for scanner.Scan() {
		idx++
		if idx == 1 {
			continue
		}
		urlChan <- strings.TrimSpace(scanner.Text())
	}

	close(urlChan)
	workerWg.Wait()
	reqWg.Wait()
}

func processIds() {
	defer workerWg.Done()

	var idx int
	for url := range urlChan {
		idx++
		if idx%1000 == 0 {
			time.Sleep(50 * time.Millisecond)
		}

		if unrestricted {
			reqWg.Add(1)
			go func(url string) {
				defer reqWg.Done()
				makeRequest(url)
			}(url)
		} else {
			makeRequest(url)
		}
	}
}

func makeRequest(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if verbose {
			fmt.Println("Could not process request:", err.Error())
		}
		return
	}

	req.Header.Add("Content-type", "application/json")

	trans := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{
		Transport: trans,
		Timeout:   timeout,
	}

	st := time.Now().UTC()
	resp, err := client.Do(req)
	atomic.AddInt64(&requestTime, time.Since(st).Nanoseconds())
	atomic.AddInt32(&reqCount, 1)

	if err != nil {
		atomic.AddInt32(&connectErrors, 1)
		if verbose {
			fmt.Println("Request Error =", err.Error())
		}
		return
	} else if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			atomic.AddInt32(&notFoundErrors, 1)
		} else {
			atomic.AddInt32(&not200Errors, 1)
		}
		if verbose {
			fmt.Println("Status =", resp.StatusCode)
		}
		return
	}

	atomic.AddInt32(&success, 1)

	// to mimick the work of processing the user
	time.Sleep(50 * time.Millisecond)
}
