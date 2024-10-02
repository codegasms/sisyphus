package main

import (
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

const MAX_WAIT int64 = 10

func main() {
	log.Println("starting load testing")

	address := os.Getenv("HOST")
	if len(address) == 0 {
		log.Fatal("HOST env not set")
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	for {
		go TestLoad(client, address)
		time.Sleep(time.Duration(rand.Int64N(MAX_WAIT)) * time.Millisecond)
	}
}

func TestLoad(client *http.Client, address string) {
	url := "http://" + address + "/"
	resp, err := client.Get(url)
	if err != nil {
		log.Println("error while GET", url)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	log.Println("GET", url, resp.Status, string(body))
}
