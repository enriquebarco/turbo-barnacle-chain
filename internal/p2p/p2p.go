package p2p

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
)

func StartServer(nodeID string, bc *blockchain.Blockchain) {
	listener, err := net.Listen("tcp", ":"+nodeID)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("Listening on %s...\n", nodeID)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, bc)
	}
}

func ConnectToNode(nodeAddress string, message string) {
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Printf("Error connecting to node: %v\n", err)
		return // Just return, don't exit
	}
	defer conn.Close()

	fmt.Fprint(conn, message)
	log.Printf("Sent message: %s to %s\n", message, nodeAddress)
}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	defer conn.Close()
	log.Printf("Connection established from %v\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("Received: %s\n", msg)

		// Handle different message types here
		// e.g., if the message is a new block, add it to the blockchain
	}
}
