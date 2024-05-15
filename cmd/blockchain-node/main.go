package main

import (
	"flag"
	"fmt"

	"github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
	"github.com/enriquebarco/turbo-barnacle-chain/internal/p2p"
)

func main() {
	// Parse command-line arguments
	localNodeName := flag.String("name", "Node", "the name you identify with")
	localPort := flag.String("port", "3000", "port to listen on")
	remoteNodeIP := flag.String("connect", "", "IP:port of remote node to connect to")
	flag.Parse()

	// Initialize the blockchain
	bc := blockchain.CreateBlockchain(2)

	// Start the P2P server on the specified port, this all deals with incoming connections
	go p2p.StartServer(*localPort, *localNodeName, *remoteNodeIP, &bc)

	// Connect to the IP addressed that was specified and say hello
	go p2p.ConnectToNode(*remoteNodeIP, *localNodeName, "REQUEST_CHAIN", fmt.Sprintf("Hello from node %s!", *remoteNodeIP))

	// handle user input to send messages to other nodes
	go p2p.HandleUserInput(&bc, *remoteNodeIP, *localNodeName)

	select {} // block main thread from closing
}
