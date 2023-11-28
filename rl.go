package main

import (
	"fmt"
	"net/http"
	"strings"
)

func HandleUnlimited(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unlimited, Baby!\n"))
	r.Body.Close()
}

func HandleLimited(w http.ResponseWriter, r *http.Request) {
	ipAddr := strings.SplitN(r.RemoteAddr, ":", 2)[0]

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
