package main

import (
	"context"
	"flag"
	"log"

	"github.com/codenotary/immudb/pkg/client"
)

var Client client.ImmuClient

func initDB(ip string, port int) {
	// even though the server address and port are defaults, setting them as a reference
	opts := client.DefaultOptions().WithAddress(ip).WithPort(port)

	c := client.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := c.OpenSession(context.Background(), []byte("immudb"), []byte("immudb"), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}
	Client = c
}

func main() {
	ipFlag := flag.String("ip", "127.0.0.1", "ip to connect to the ImmuDB instance")
	portFlag := flag.Int("port", 3322, "port to connect to the ImmuDB instance")
	topicFlag := flag.String("net", "ImmuDBTopic", "name of the topic for the network")
	flag.Parse()

	NET = *topicFlag
	GetNode()

	initDB(*ipFlag, *portFlag)

	// ensure connection is closed
	defer Client.CloseSession(context.Background())

	// PromptUI to select action
	TopUI()
}
