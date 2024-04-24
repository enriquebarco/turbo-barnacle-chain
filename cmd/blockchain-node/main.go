package main

import (
	// "github.com/enriquebarco/turbo-barnacle-chain/internal/blockchain"
	"github.com/enriquebarco/turbo-barnacle-chain/internal/p2p"
)

func main() {

	// Initialize the blockchain
	// bc := blockchain.CreateBlockchain(2)

	// Start the P2P server on the specified port
	p2p.StartServer()
}
