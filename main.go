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
	buckets = make(map[string]uint32)
	mutex = sync.RWMutex{}

	slidingMap = make(map[string][]time.Time)
	windowCount = make(map[string]uint32)
	previousWindow = make(map[string]uint32)
	lastUpdate = time.Now()
)

func main() {
	server := http.NewServeMux()
	server.Handle("/unlimited", http.HandlerFunc(HandleUnlimited))
	server.Handle("/limited", http.HandlerFunc(HandleLimited))

	// go addTokens()
	// go updateFixedWindow()
	// go updateSlidingWindow()
	go updateSlidingCounter()

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

func updateFixedWindow() {
	for {
		mutex.Lock()
		windowCount = make(map[string]uint32)
		mutex.Unlock()

		time.Sleep(WINDOW)
	}
}

func updateSlidingWindow() {
	for {
		mutex.Lock()
		for k, v := range slidingMap {
			times := v
			// var cutoff int
			for i := range times {
				if times[i].Before(time.Now().Add(-WINDOW)) {
					continue
				}
				slidingMap[k] = slidingMap[k][i:]
				break
			}
			// slidingMap[k] = times[cutoff:]
		}
		mutex.Unlock()

		time.Sleep(waitDuration)
	}
}

func updateSlidingCounter() {
	for {
		mutex.Lock()
		previousWindow = windowCount
		windowCount = make(map[string]uint32)
		lastUpdate = time.Now()
		mutex.Unlock()

		time.Sleep(WINDOW)
	}
}
