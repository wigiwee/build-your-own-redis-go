package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	//reading cmd args
	port := flag.Int("port", 6379, "Redis server port")
	flag.Parse()

	fmt.Println(*port)
	store := storage.NewStore()
	handler := command.NewHandler(store)
	srv := server.NewServer(*handler)

	if err := srv.Start(":" + strconv.Itoa(*port)); err != nil {
		log.Fatal(err)
	}
}
