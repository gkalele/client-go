package main

import (
	"fmt"
	"net/http"
)

func StartWebServer(podName string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server - podname %s", podName)
	})
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Could not start HTTP server - %s", err.Error())
	}
}
