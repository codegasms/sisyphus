package main

import (
	"io"
	"log"
	"net"
	"os"
)

func getEnvOr(key, or string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	} else {
		return or
	}
}

func main() {
	config, err := LoadConfig("data/config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded config", config)

	strategy, err := StrategyFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	serverAddr, err := strategy.ServerAddr()
	log.Println(serverAddr)

	port := getEnvOr("PORT", "3030")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("TCP listener initialized at", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func(conn net.Conn) {
			_, err := io.Copy(conn, conn)
			if err != nil {
				log.Println(err)
			}
			// Close the connection and decrement connection count.
			conn.Close()
		}(conn)
	}
}
