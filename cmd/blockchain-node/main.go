package main

import (
	"flag"
	"fmt"

	"github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
	"github.com/enriquebarco/turbo-barnacle-chain/internal/p2p"
)

func main() {
	// Parse command-line arguments
	localPort := flag.String("port", "3000", "port to listen on")
	remoteNode := flag.String("connect", "", "IP:port of remote node to connect to")
	flag.Parse()

	// Initialize the blockchain
	bc := blockchain.CreateBlockchain(2)

	// Start the P2P server on the specified port
	go p2p.StartServer(*localPort, &bc)

	// If a remote node address is provided, connect to it
	if *remoteNode != "" {
		go p2p.ConnectToNode(*remoteNode, fmt.Sprintf("Hello from node %s!", *localPort))
	}

	// Start handling user input to send messages to other nodes
	p2p.HandleUserInput(&bc, *remoteNode)
}
