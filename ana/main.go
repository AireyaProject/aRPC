package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const SockAddr = "/tmp/aireya.sock"

// RouteRequest is what the C++ SDK sends to the Agent
type RouteRequest struct {
	Service string `json:"service"`
}

// RouteResponse is what the Agent replies to the SDK
type RouteResponse struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

// Mock Database of Global Routing Rules
// In a real agent, this maps to Etcd or an xDS control plane stream.
var routingTable = map[string]RouteResponse{
	"UserService": {Host: "127.0.0.1", Port: 8080},
	"OrderService": {Host: "10.0.0.5", Port: 9090},
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	log.Printf("Aireya Node Agent (ANA) started on %s\n", SockAddr)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigc
		os.RemoveAll(SockAddr)
		os.Exit(0)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req RouteRequest
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Decode error: %v\n", err)
		return
	}

	log.Printf("SDK requested route for service: %s\n", req.Service)

	resp, ok := routingTable[req.Service]
	if !ok {
		resp = RouteResponse{Host: "", Port: 0} // Not found
	}

	if err := encoder.Encode(resp); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}
