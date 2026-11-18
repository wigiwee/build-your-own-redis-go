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

	serverConfig := config.ServerConfig{}
	//reading cmd args
	port := flag.Int("port", 6379, "Redis server port")
	masterAddr := flag.String("replicaof", "", "Redis server role")

	flag.Parse()

	serverConfig.Port = *port
	if *masterAddr == "" {
		serverConfig.ReplicationConfig.Role = "master"
		serverConfig.ReplicationConfig.Master_replid = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		serverConfig.ReplicationConfig.Master_repl_offset = 0
	} else {
		serverConfig.ReplicationConfig.Role = "slave"
	}

	store := storage.NewStore()
	handler := command.NewHandler(store, serverConfig)
	srv := server.NewServer(*handler)

	if err := srv.Start(*port); err != nil {
		log.Fatal(err)
	}
}
