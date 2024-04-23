package p2p

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
)

// StartServer starts the P2P server and listens for incoming connections.
func StartServer(nodeID string, bc *blockchain.Blockchain) {
	ln, err := net.Listen("tcp", ":"+nodeID)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Listening for P2P connections on %s...\n", nodeID)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn, bc)
	}
}

// handleConnection deals with incoming data.
func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	defer conn.Close()

	address := conn.RemoteAddr().String()
	log.Printf("New connection established from %s\n", address)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Text()
		log.Printf("Received from %s: %s\n", address, data)

		// Process incoming data here
		// ...
	}
}

// ConnectToNode connects to a specified node and sends a message.
func ConnectToNode(nodeAddress, message string) {
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Printf("Error connecting to node at %s: %v\n", nodeAddress, err)
		return
	}
	defer conn.Close()

	fmt.Fprintln(conn, message)
	log.Printf("Sent message to %s: %s\n", nodeAddress, message)
}

// HandleUserInput allows the user to input messages to be sent to the network.
func HandleUserInput(bc *blockchain.Blockchain, nodeAddress string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read user input: %v\n", err)
			continue
		}

		message := strings.TrimSpace(input)
		ConnectToNode(nodeAddress, message)
	}
}
