package main

import (
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	srv := server.NewServer()
	if err := srv.Start(":6379"); err != nil {
		log.Fatal(err)
	}

}
