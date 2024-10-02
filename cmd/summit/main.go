package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ApiResponse struct {
	Message string `json:"message"`
}

func LogRequest(req *http.Request) {
	log.Printf("%v -> %v %v %v", req.RemoteAddr, req.Method, req.URL, req.Proto)
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		LogRequest(req)

		res := ApiResponse{"ok"}
		resJson, _ := json.Marshal(res)
		w.Write(resJson)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		LogRequest(req)

		res := ApiResponse{"running"}
		resJson, _ := json.Marshal(res)
		w.Write(resJson)
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT env not set")
	}

	log.Printf("HTTP server listening on localhost:%v", port)
	http.ListenAndServe(":"+port, nil)
}
