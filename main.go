package main

import (
	"fmt"
	"net/http"
)

var (
	PORT = "8080"
)

func main() {
	server := http.NewServeMux()
	server.Handle("/unlimited", http.HandlerFunc(HandleUnlimited))
	server.Handle("/limited", http.HandlerFunc(HandleLimited))

	fmt.Println("Listening on port " + PORT)
	err := http.ListenAndServe("localhost:" + PORT, server)
	if err != nil {
		fmt.Println(err.Error())
	}
}
