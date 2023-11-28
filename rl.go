package main

import (
	"net/http"
)


func HandleUnlimited(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unlimited, Baby!\n"))
	r.Body.Close()
}

func HandleLimited(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Limited, unfortunately\n"))
	r.Body.Close()
}
