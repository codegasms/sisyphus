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

	port := getEnvOr("PORT", "8000")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("TCP listener initialized at", listener.Addr())

	BalanceLoad(listener, strategy)
}

func BalanceLoad(listener net.Listener, strategy Strategy) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go BalanceOneLoad(conn, strategy)
	}
}

func BalanceOneLoad(conn1 net.Conn, strategy Strategy) {
	defer conn1.Close()

	serverAddr, err := strategy.ServerAddr(conn1.RemoteAddr().String())
	if err != nil {
		log.Println("no servers to forward to:", err)
		return
	}

	conn2, err := net.Dial("tcp", string(serverAddr))
	if err != nil {
		log.Println("couldn't connect to the selected server", serverAddr)
		return
	}
	defer conn2.Close()

	log.Printf("forwarding %v to %v", conn1.RemoteAddr(), serverAddr)

	// Increment and defer decrement the connection count.
	strategy.Connected(serverAddr)
	defer strategy.Disconnected(serverAddr)

	go io.Copy(conn1, conn2)
	io.Copy(conn2, conn1)

	log.Printf("connection %v to %v closed", conn1.RemoteAddr(), serverAddr)
}
