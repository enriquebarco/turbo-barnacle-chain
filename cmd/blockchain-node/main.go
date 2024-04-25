package main

import (
	"flag"
	"fmt"

	"github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
	"github.com/enriquebarco/turbo-barnacle-chain/internal/p2p"
)

func main() {
	// Parse command-line arguments
	nodeName := flag.String("name", "Node", "the name you identify with")
	localPort := flag.String("port", "3000", "port to listen on")
	remoteNodeIP := flag.String("connect", "", "IP:port of remote node to connect to")
	flag.Parse()

	// Initialize the blockchain
	bc := blockchain.CreateBlockchain(2)

	// Start the P2P server on the specified port
	go p2p.StartServer(*localPort, *nodeName, &bc)

	// If a remote node address is provided, connect to it
	if *remoteNodeIP != "" {
		go p2p.ConnectToNode(*remoteNodeIP, *nodeName, fmt.Sprintf("Hello from node %s!", *remoteNodeIP))
	}

	// Start handling user input to send messages to other nodes
	p2p.HandleUserInput(&bc, *remoteNodeIP, *nodeName)
}
