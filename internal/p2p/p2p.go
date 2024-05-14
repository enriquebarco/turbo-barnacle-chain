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

var userCommands = map[string]Command{
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
			// Format the message as a transaction
			message := fmt.Sprintf("TXN:%s,%s,%f", args[0], args[1], amount)
			if conn != nil {
				fmt.Fprintf(conn, "%s\n", message)
			}
		},
	},
	// add more commands as needed
}

// StartServer starts the P2P server and listens for incoming connections.
func StartServer(nodeID string, nodeName string, bc *blockchain.Blockchain) {
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

		// we create a new go routine so that we can handle multiple connections simultaneously, each independent from the others.
		go handleConnection(conn, bc)
	}
}

// handleConnection deals with incoming data.
func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	defer conn.Close()
	// log.Printf("New connection established from %s\n", nodeName)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// full message
		fullMessage := scanner.Text()
		messageParts := strings.SplitN(fullMessage, ":", 2) // Split to get the sender's name and the message

		if len(messageParts) != 2 {
			log.Printf("Invalid message format received: %s\n", fullMessage)
			continue
		}

		senderName := messageParts[0]
		message := messageParts[1]

		fmt.Printf("\033[32m%s: %s\033[0m\n", senderName, message)
		// detect if a user has sent a valid txn, add it to the blockchain
		if strings.HasPrefix(message, "TXN:") {
			// check that a valid transaction happened on the chain by checking the chain is still valid
			// print the entire block chain
			parts := strings.Split(message[4:], ",")
			if len(parts) == 3 {
				from := parts[0]
				to := parts[1]
				amount, err := strconv.ParseFloat(parts[2], 64)
				if err == nil {
					bc.AddBlock(from, to, amount)
					fmt.Println("New block received and added to the blockchain")
					fmt.Printf("From: %s, To: %s, Amount: %f\n", from, to, amount)
					log.Printf("Block added for transaction from %s to %s of %f", from, to, amount)

					fmt.Println("Current Blockchain:")
					for i, block := range bc.Chain {
						if i == 0 {
							continue // Skip the genesis block
						}
						fmt.Printf("Data: %v\n", block.Data)
					}
					continue
				}
			}
			log.Println("Invalid transaction format")
		}
	}
}

// ConnectToNode connects to a specified node and sends a message
func ConnectToNode(nodeAddress, nodeName, message string) {
	// establish a tcp connection
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Printf("Error connecting to node at %s: %v\n", nodeAddress, err)
		return
	}
	defer conn.Close()

	formattedMessage := fmt.Sprintf("%s:%s", nodeName, message)
	fmt.Fprintf(conn, "%s\n", formattedMessage)
}

// HandleUserInput allows the user to input messages to be sent to the network.
func HandleUserInput(bc *blockchain.Blockchain, nodeAddress string, nodeName string) {
	// create a buffered reader that reads from the standard input
	reader := bufio.NewReader(os.Stdin)

	// run an infinite for loop that continous reading input from the user until the program is terminated
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
			// the expected format is -- send from,to,amount
			details := message[5:]
			parts := strings.Split(details, ",")
			if len(parts) == 3 {
				from := parts[0]
				to := parts[1]
				amount, err := strconv.ParseFloat(parts[2], 64)
				if err != nil {
					fmt.Println("Invalid amount:", parts[2])
					continue
				}

				// Add the block to the local blockchain
				bc.AddBlock(from, to, amount)
				fmt.Println("New block added to the blockchain")
				fmt.Printf("From: %s, To: %s, Amount: %f\n", from, to, amount)
				fmt.Println("Current Blockchain:")
				for i, block := range bc.Chain {
					if i == 0 {
						continue // Skip the genesis block
					}
					fmt.Printf("Data: %v\n", block.Data)
				}

				// Broadcast the new block to the connected node
				if nodeAddress != "" {
					message := fmt.Sprintf("TXN:%s,%s,%f", from, to, amount)
					ConnectToNode(nodeAddress, nodeName, message)
				}
				continue
			}
			fmt.Println("Invalid command format. Use: send from,to,amount (e.g. send mel,kike,10)")
		}
		// if not a specific blockchain message, this handles it
		if nodeAddress != "" {
			ConnectToNode(nodeAddress, nodeName, message)
		} else {
			fmt.Println("No remote node address specified. Unable to send the message.")
		}
	}
}
