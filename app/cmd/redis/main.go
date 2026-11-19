package main

import (
	"flag"
	"log"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

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
		split := strings.Split(*masterAddr, "")
		serverConfig.ReplicationConfig.Master_host = split[0]
		serverConfig.ReplicationConfig.Master_port = split[1]

	}
	store := storage.NewStore()
	handler := command.NewHandler(store, serverConfig)
	srv := server.NewServer(*handler)

	if err := srv.Start(*port); err != nil {
		log.Fatal(err)
	}
}
