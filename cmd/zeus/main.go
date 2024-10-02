package main

import (
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("starting load testing")

	address := os.Getenv("HOST")
	if len(address) == 0 {
		log.Fatal("HOST env not set")
	}

	for {
		go TestLoad(address)
		time.Sleep(time.Duration(rand.Int64N(50)) * time.Millisecond)
	}
}

func TestLoad(address string) {
	route := "http://" + address + "/"
	res, err := http.Get(route)
	if err != nil {
		log.Println("error while GET", route)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	log.Println("GET", route, string(body))
}
