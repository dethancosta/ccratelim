package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	waitDuration = 1 * time.Second
	WINDOW       = 60 * time.Second
	WINDOW_LIMIT = 60
	REDIS_ADDR   = "localhost:6379"
)

var (
	PORT    = "8080"
	buckets = make(map[string]uint32)
	mutex   = sync.RWMutex{}
	ctx     = context.Background()

	slidingMap     = make(map[string][]time.Time)
	windowCount    = make(map[string]uint32)
	previousWindow = make(map[string]uint32)
	redisClient    = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: "",
		DB:       0,
	})
)

func main() {
	flag.StringVar(&PORT, "p", "8080", "port that the server will run on")
	flag.Parse()
	server := http.NewServeMux()
	server.Handle("/unlimited", http.HandlerFunc(HandleUnlimited))
	server.Handle("/limited", http.HandlerFunc(HandleLimited))

	if redisClient.Echo(ctx, "Hi").Val() != "Hi" {
		fmt.Println("Echo failed")
	} else {
		fmt.Println("Echo successful")
	}
	redisClient.HSet(ctx, "prevWindow", "") // TODO check return value
	redisClient.HSet(ctx, "currWindow", "") // TODO check return value
	redisClient.Set(ctx, "lastUpdate", time.Now(), 0)

	// go addTokens()
	// go updateFixedWindow()
	// go updateSlidingWindow()
	go updateSlidingCounter()

	fmt.Println("Listening on port " + PORT)
	err := http.ListenAndServe("localhost:"+PORT, server)
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
		time.Sleep(WINDOW)
		currentWindow := redisClient.HGetAll(ctx, "currWindow")
		redisClient.HSet(ctx, "prevWindow", currentWindow) // TODO check return value
		redisClient.HSet(ctx, "currWindow", "")            // TODO check return value
		redisClient.Set(ctx, "lastUpdate", time.Now(), 0)

		/*
			previousWindow = windowCount
			windowCount = make(map[string]uint32)
		*/
	}
}
