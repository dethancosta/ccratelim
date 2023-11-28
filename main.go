package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	waitDuration = 1 * time.Second
	PORT = "8080"
	WINDOW = 60 * time.Second
	WINDOW_LIMIT = 60
)

var (
	buckets = make(map[string]uint8)
	mutex = sync.RWMutex{}

	windowSlice = make([]RequestItem, 0)
	windowCount = make(map[string]uint32)
)

func main() {
	server := http.NewServeMux()
	server.Handle("/unlimited", http.HandlerFunc(HandleUnlimited))
	server.Handle("/limited", http.HandlerFunc(HandleLimited))

	//go addTokens()
	go updateWindow()

	fmt.Println("Listening on port " + PORT)
	err := http.ListenAndServe("localhost:" + PORT, server)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func addTokens() {
	for {
	mutex.Lock()
	for i := range buckets {
		if buckets[i] < 10 {
			buckets[i]++
		}
	}
	mutex.Unlock()
	time.Sleep(waitDuration)
	}
}

func updateWindow() {
	for {
		stale := make([]RequestItem, 0)
		mutex.Lock()
		for i := range windowSlice {
			if windowSlice[i].At.Compare(time.Now().Add(-WINDOW)) < 0 {
				continue 
			}
			stale = windowSlice[:i]
			windowSlice = windowSlice[i:]
		}
		for _, r := range stale {
			windowCount[r.IpAddr]--
		}
		mutex.Unlock()

		time.Sleep(waitDuration)
	}
}

