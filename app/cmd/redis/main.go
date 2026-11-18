package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
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

	store := storage.NewStore()
	handler := command.NewHandler(store, config.ServerConfig{
		Port:        *port,
		Version:     "0.1.0",
		Mode:        "standalone",
		BindAddress: "127.0.0.1",
		ReplicationConfig: config.ReplicationConfig{
			Role: "master",
		},
	})
	srv := server.NewServer(*handler)

	if err := srv.Start(*port); err != nil {
		log.Fatal(err)
	}
}
