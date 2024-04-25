package p2p

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
)

type TransactionMessage struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type Command struct {
	Execute func(args []string, bc *blockchain.Blockchain, conn net.Conn)
}

var commands = map[string]Command{
	"send": {
		Execute: func(args []string, bc *blockchain.Blockchain, conn net.Conn) {
			if len(args) != 3 {
				fmt.Println("Usage: send <from> <to> <amount>")
				return
			}
			amount, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Println("Invalid amount:", args[2])
				return
			}
			bc.AddBlock(args[0], args[1], amount)
			if conn != nil {
				fmt.Fprintf(conn, "TXN:%s,%s,%f\n", args[0], args[1], amount)
			}
		},
	},
	// Add more commands as needed
}

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
		message := scanner.Text()
		log.Printf("Received from %s: %s\n", address, message)

		// detect if a user has sent a valid txn, add it to the blockchain
		if strings.HasPrefix(message, "send:") {
			// check that a valid transaction happened on the chain by checking the chain is still valid
			// print the entire block chain
			parts := strings.Split(message[4:], ",")
			if len(parts) == 3 {
				from := parts[0]
				to := parts[1]
				amount, err := strconv.ParseFloat(parts[2], 64)
				if err == nil {
					bc.AddBlock(from, to, amount)
					log.Printf("Block added for transaction from %s to %s of %f", from, to, amount)
					continue
				}
			}
			log.Println("Invalid transaction format")
		}
	}
}

// ConnectToNode connects to a specified node and sends a message
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
		// handle a user specifically trying to send a transaction on the blockchain
		if strings.HasPrefix(message, "send") {
			// the expected format is send from,to,amount
			details := message[5:]
			parts := strings.Split(details, ",")
			if len(parts) == 3 {
				if nodeAddress != "" {
					ConnectToNode(nodeAddress, "TXN:"+details)
				}
				continue
			}
			fmt.Println("Invalid command format. Use: send from,to,amount (e.g. send from mel,kike,10)")
		}
		if nodeAddress != "" {
			ConnectToNode(nodeAddress, message)
		} else {
			fmt.Println("No remote node address specified. Unable to send the message.")
		}
	}
}
