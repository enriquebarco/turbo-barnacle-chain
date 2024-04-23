package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

func (b Block) calculateHash() string {
	// convert block data into json
	data, _ := json.Marshal(b.data)
	// Concatenate the previous block’s hash, and the current block’s data, timestamp, and PoW
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	// hash with sha256 algo
	blockHash := sha256.Sum256([]byte(blockData))
	// return the base 16 hash as a string
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	// continue to change the proof of work value of the current block until we satisfiy our mining condition of (starting zeros > difficulty)
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

// create the genesis block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0 Hello Mel",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

// adding new blocks to the blockchain
func (b *Blockchain) addBlock(from, to string, amount float64) {
	// collect details of a transaction (sender, receiver, and transfer amount)
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	// create a new block with the transaction details
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	// mine the new block with the previous block hash, current block data, and generated PoW
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

// check the validity of the blockchain. No transactions should be tampered with
func (b *Blockchain) isValid() bool {
	// skip genesis block because it does not have a previous block
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		// first, recalculate the hash of the block and compare it to the stored hash value
		// second, check if the previous hash value saved in this block is equal to its previous hash
		// if a block has been tampered with, this check willf fail since the hash will change
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

// check the validity of a new block
func IsValidNewBlock(newBlock, previousBlock Block) bool {
	// Check if the previous hash matches
	if previousBlock.hash != newBlock.previousHash {
		return false
	}

	// Check if the hash of the new block is correct
	if newBlock.hash != newBlock.calculateHash() {
		return false
	}

	return true
}

// TODO: create a peer 2 peer network to connect nodes
// TODO: implement a consensus algorithm
// TODO: create security countermeasures
