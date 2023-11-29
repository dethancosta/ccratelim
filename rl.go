package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func HandleUnlimited(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unlimited, Baby!\n"))
	r.Body.Close()
}

func HandleLimited(w http.ResponseWriter, r *http.Request) {
	ipAddr := strings.SplitN(r.RemoteAddr, ":", 2)[0]

	// TokenBucket(ipAddr, w)
	// FixedWindow(ipAddr, w)
	SlidingWindow(ipAddr, w)
}

type RequestItem struct {
	IpAddr string
	At time.Time
}


// Rate Limiter functions

func TokenBucket(ipAddr string, w http.ResponseWriter) {
	mutex.Lock()
	if _, ok := buckets[ipAddr]; !ok {
		buckets[ipAddr] = 10
	}

	count := buckets[ipAddr]
	if count > 0 {
		buckets[ipAddr] -= 1
	}
	mutex.Unlock()
	if count > 0 {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("%d tokens", count)))
		return
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
}

func FixedWindow(ipAddr string, w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()
	count, ok := windowCount[ipAddr]
	if !ok {
		windowCount[ipAddr] = 0
	} else if count >= WINDOW_LIMIT {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	windowCount[ipAddr]++
	fmt.Printf("%s: %d\n", ipAddr, count+1)
	w.WriteHeader(http.StatusOK)
}

func SlidingWindow(ipAddr string, w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()
	l := len(slidingMap[ipAddr])
	if l >= WINDOW_LIMIT {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	fmt.Printf("%s: %d\n", time.Now().Format(time.TimeOnly), len(slidingMap[ipAddr]))
	slidingMap[ipAddr] = append(slidingMap[ipAddr], time.Now())
	w.WriteHeader(http.StatusOK)
}
