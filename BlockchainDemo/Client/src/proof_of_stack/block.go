package proof_of_stack

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Validator string
	Data      []byte
}

// SHA256 hasing
// calculateHash is a simple SHA256 hashing function
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//calculateBlockHash returns the hash of all block information
func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	return calculateHash(record)
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, BPM int, address string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Validator = address

	newBlock.Data = []byte("Test data")
	newBlock.Hash = calculateBlockHash(newBlock)

	return newBlock, nil
}

// isBlockValid makes sure block is valid by checking index
// and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {

	if oldBlock.Index+1 != newBlock.Index {
		fmt.Println("Invalid index !!", newBlock.Index, oldBlock.Index)
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		fmt.Println("Invalid prev  hash !!")
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		fmt.Println("Invalid  hash !!")
		return false
	}
	return true
}

func (b *Block) NewInstance() Block {
	return Block{}
}
