package p2p

import (
	"bufio"
	"encoding/json"
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

// StartServer starts the P2P server and listens for incoming connections.
func StartServer(localPort, nodeName, remoteNodeIP string, bc *blockchain.Blockchain) {
	ln, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Listening for P2P connections on %s...\n", localPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		// we create a new go routine so that we can handle multiple connections simultaneously, each independent from the others.
		go handleConnection(conn, bc, nodeName, remoteNodeIP)
	}
}

// handleConnection deals with incoming data.
func handleConnection(conn net.Conn, bc *blockchain.Blockchain, nodeName, remoteNodeIP string) {
	defer conn.Close()
	// log.Printf("New connection established from %s\n", nodeName)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// full message
		fullMessage := scanner.Text()
		messageParts := strings.SplitN(fullMessage, ":", 3) // Split to get the sender's name, optional message type, and the message

		if len(messageParts) < 2 {
			log.Printf("Invalid message format received: %s\n", fullMessage)
			continue
		}

		senderName := messageParts[0]
		message := messageParts[1]
		var messageType string
		if len(messageParts) == 3 {
			messageType = messageParts[1]
			message = messageParts[2]
		} else {
			messageType = "MESSAGE"
		}

		fmt.Printf("\033[32m%s [%s]: %s\033[0m\n", senderName, messageType, message)

		switch messageType {
		case "REQUEST_CHAIN":
			fmt.Println("Remote node has requested our blockchain")
			// send the blockchain to the remote node
			chainJson, err := json.Marshal(bc.Chain)
			if err != nil {
				log.Printf("Failed to marshal blockchain: %v\n", err)
				return
			}
			fmt.Println("Sending blockchain to remote node...")
			if err := ConnectToNode(remoteNodeIP, nodeName, "RECEIVE_CHAIN", string(chainJson)); err != nil {
				log.Printf("Blockchain rejected: %v\n", err)
			}

		case "RECEIVE_CHAIN":
			fmt.Println("Received blockchain from remote node")
			var newChain []blockchain.Block
			err := json.Unmarshal([]byte(message), &newChain)
			if err != nil {
				log.Printf("Failed to unmarshal blockchain: %v\n", err)
				return
			}
			err = bc.ReplaceChain(newChain)
			if err != nil {
				log.Printf("Failed to replace chain: %v\n", err)
			} else {
				fmt.Println("Blockchain replaced with the received chain")
				fmt.Println("Current Blockchain:")
				bc.PrintChain()
			}

		case "RECIEVE_BLOCK":
			var newBlock blockchain.Block
			err := json.Unmarshal([]byte(message), &newBlock)
			if err != nil {
				log.Printf("Failed to unmarshal block: %v\n", err)
				return
			}
			err = bc.ReceiveBlock(newBlock)
			if err != nil {
				log.Printf("Failed to receive block: %v\n", err)
			} else {
				fmt.Println("New block added to the blockchain")
				fmt.Printf("From: %s, To: %s, Amount: %f\n", newBlock.Data["from"], newBlock.Data["to"], newBlock.Data["amount"])
				fmt.Println("Current Blockchain:")
				bc.PrintChain()
			}
		case "MESSAGE":
			// fmt.Printf("\033[32m%s: %s\033[0m\n", senderName, message)
		default:
			fmt.Printf("Unknown message type: %s\n", messageType)
		}
	}
}

// ConnectToNode connects to a specified node and sends a message
func ConnectToNode(nodeAddress, nodeName, messageType string, message string) error {
	// establish a tcp connection
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Printf("Error connecting to node at %s: %v\n", nodeAddress, err)
		return err
	}
	defer conn.Close()

	var formattedMessage string
	if messageType == "" {
		formattedMessage = fmt.Sprintf("%s:%s", nodeName, message)
	} else {
		formattedMessage = fmt.Sprintf("%s:%s:%s", nodeName, messageType, message)
	}

	_, err = fmt.Fprintf(conn, "%s\n", formattedMessage)
	if err != nil {
		log.Printf("Error sending message to node at %s: %v\n", nodeAddress, err)
		return err
	}

	return nil
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
				fmt.Println("Current Blockchain:")
				bc.PrintChain()
				// get the new block and broadcast it to the connected node
				lastBlock := bc.Chain[len(bc.Chain)-1]

				blockJson, err := json.Marshal(lastBlock)
				if err != nil {
					fmt.Println("Error marshalling block to JSON:", err)
					return
				}

				if nodeAddress != "" {
					fmt.Println("Broadcasting new block to the network...")
					message := string(blockJson)
					if err := ConnectToNode(nodeAddress, nodeName, "RECIEVE_BLOCK", message); err != nil {
						fmt.Println("Failed to broadcast block:", err)
					}
				}
				continue
			}
			fmt.Println("Invalid command format. Use: send from,to,amount (e.g. send mel,kike,10)")
		}
		// handle requesting the blockchain
		if strings.HasPrefix(message, "request_chain") {
			fmt.Println("Requesting latest blockchain from network...")
			if nodeAddress != "" {
				if err := ConnectToNode(nodeAddress, nodeName, "REQUEST_CHAIN", ""); err != nil {
					fmt.Println("Failed to request blockchain:", err)
				}
			} else {
				fmt.Println("No remote node address specified. Unable to request the blockchain.")
			}
		}
		// if not a specific blockchain message, this handles it
		if nodeAddress != "" {
			if err := ConnectToNode(nodeAddress, nodeName, "MESSAGE", message); err != nil {
				fmt.Println("Failed to send message:", err)
			}
		} else {
			fmt.Println("No remote node address specified. Unable to send the message.")
		}
	}
}
