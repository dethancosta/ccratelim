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
	// SlidingWindow(ipAddr, w)
	SlidingCount(ipAddr, w)
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
	slidingMap[ipAddr] = append(slidingMap[ipAddr], time.Now())
	w.WriteHeader(http.StatusOK)
}

func SlidingCount(ipAddr string, w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()
	ratioCurrent := time.Now().Sub(lastUpdate).Seconds()/WINDOW.Seconds()
	ratioPrev := 1.0 - ratioCurrent
	var prevCount uint32
	var currCount uint32
	var ok bool
	if prevCount, ok = previousWindow[ipAddr]; !ok {
		previousWindow[ipAddr] = 0
	}
	if currCount, ok = windowCount[ipAddr]; !ok {
		windowCount[ipAddr] = 0
	}
	rate := (ratioPrev * float64(prevCount)) +
		(ratioCurrent * float64(currCount))
	if rate > WINDOW_LIMIT {
		// http.Error(w, "Slow Down!!", http.StatusTooManyRequests)
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	windowCount[ipAddr]++
	w.WriteHeader(http.StatusOK)
}
