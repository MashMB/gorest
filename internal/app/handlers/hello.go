package handlers

import "net/http"

func Hello(res http.ResponseWriter, req *http.Request) {
	message := "Hello World"
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(message))
}
