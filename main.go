package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	PORT = "8080"
	buckets map[string]uint8 = make(map[string]uint8)
	mutex = sync.RWMutex{}
	waitDuration = 1 * time.Second
)

func main() {
	server := http.NewServeMux()
	server.Handle("/unlimited", http.HandlerFunc(HandleUnlimited))
	server.Handle("/limited", http.HandlerFunc(HandleLimited))

	go addTokens()

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
