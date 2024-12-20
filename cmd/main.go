package main

import (
	"calc-service/internal/server"
	"net/http"
)

func main () {
	http.HandleFunc("/api/v1/calculate", server.CalculateHandler)
	http.ListenAndServe(":8080", nil)
}