package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
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
		split := strings.Split(*masterAddr, " ")
		serverConfig.ReplicationConfig.Master_host = split[0]
		serverConfig.ReplicationConfig.Master_port = split[1]
		fmt.Println(split)
		conn, err := net.Dial("tcp", ":"+serverConfig.ReplicationConfig.Master_port)
		if err != nil {
			log.Println("error dialing error", err)
		}
		// reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
		reader := bufio.NewReader(conn)

		pingStr, err := utils.Serialize([]string{"PING"})
		if err != nil {
			log.Println("error serializing", err)
		}
		writer.Write([]byte(pingStr))
		writer.Flush()
		log.Println(reader.ReadLine())

		replconfStr, err := utils.Serialize([]string{"REPLCONF", "listening-port", strconv.Itoa(serverConfig.Port)})
		if err != nil {
			log.Println(err)
		}
		writer.Write([]byte(replconfStr))
		writer.Flush()
		log.Println(reader.ReadLine())

		replconfStr, err = utils.Serialize([]string{"REPLCONF", "capa", "psync2"})
		if err != nil {
			log.Println(err)
		}
		writer.Write([]byte(replconfStr))
		writer.Flush()
		log.Println(reader.ReadLine())

		replconfStr, err = utils.Serialize([]string{"PSYNC", "?", "-1"})
		if err != nil {
			log.Println(err)
		}
		writer.Write([]byte(replconfStr))
		writer.Flush()
		log.Println(reader.ReadLine())

	}
	store := storage.NewStore()
	handler := command.NewHandler(store, serverConfig)
	srv := server.NewServer(*handler)

	if err := srv.Start(*port); err != nil {
		log.Fatal(err)
	}
}
