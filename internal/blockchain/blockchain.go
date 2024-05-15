package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Nonce        int
}

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func (b Block) calculateHash() string {
	// convert block data into json
	data, _ := json.Marshal(b.Data)
	// Concatenate the previous block’s hash, and the current block’s data, timestamp, and nonce
	blockData := b.PreviousHash + string(data) + b.Timestamp.UTC().String() + strconv.Itoa(b.Nonce)
	// hash with sha256 algo
	blockHash := sha256.Sum256([]byte(blockData))
	// return the base 16 hash as a string
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	// continue to change the proof of work value of the current block until we satisfiy our mining condition of (starting zeros > difficulty)
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}
}

// create the genesis block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		Hash:      "0 Hello Mel",
		Timestamp: time.Now().UTC(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

// adding new blocks to the blockchain
func (b *Blockchain) AddBlock(from, to string, amount float64) {
	// collect details of a transaction (sender, receiver, and transfer amount)
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	// create a new block with the transaction details
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         blockData,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().UTC(),
	}
	// mine the new block with the previous block hash, current block data, and generated nonce
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

// recieiving blocks from other nodes
func (bc *Blockchain) ReceiveBlock(newBlock Block) error {
	err := isValidNewBlock(newBlock, bc.Chain[len(bc.Chain)-1])
	if err != nil {
		return fmt.Errorf("failed to validate new block: %w", err)
	}
	bc.Chain = append(bc.Chain, newBlock)
	return nil
}

// check the validity of the blockchain. No transactions should be tampered with
func isValidChain(chain []Block) bool {
	// skip genesis block because it does not have a previous block
	for i := range chain[1:] {
		previousBlock := chain[i]
		currentBlock := chain[i+1]
		// first, recalculate the hash of the block and compare it to the stored hash value
		// second, check if the previous hash value saved in this block is equal to its previous hash
		// if a block has been tampered with, this check willf fail since the hash will change
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// check the validity of a new block
func isValidNewBlock(newBlock, previousBlock Block) error {
	// Check if the previous hash matches
	if previousBlock.Hash != newBlock.PreviousHash {
		return errors.New("previous hash does not match")
	}

	// Check if the hash of the new block is correct
	if newBlock.Hash != newBlock.calculateHash() {
		return errors.New("block hash is invalid")
	}

	return nil
}

// print the chain to the console
func (bc *Blockchain) PrintChain() {
	for i, block := range bc.Chain {
		if i == 0 {
			continue // Skip the genesis block
		}
		fmt.Printf("Transaction: %v, Nonce: %d\n", block.Data, block.Nonce)
	}
}

func (bc *Blockchain) ReplaceChain(newChain []Block) error {
	// Check if the new chain is longer than the current chain
	if len(newChain) > len(bc.Chain) {
		// Check if the new chain is valid
		if !isValidChain(newChain) {
			return errors.New("received invalid chain")
		}
		// Replace the current chain with the new chain
		bc.Chain = newChain
	} else {
		return errors.New("blockchain up to date, chains are same length")
	}
	return nil
}
