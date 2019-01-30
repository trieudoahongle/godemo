package proof_of_stack

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"
	"util"
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

	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Data = makeSignature(newBlock.Hash)

	return newBlock, nil
}
func makeSignature(hashString string) []byte {
	hashed, _ := hex.DecodeString(hashString)
	rng := rand.Reader
	priv := util.GetPrivateKey("serverpriv2048.txt")
	//	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rng, priv, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return nil
	}

	return signature

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
	if !isServerData(newBlock.Data, newBlock.Hash) {
		return false
	}
	return true
}

func isServerData(data []byte, hashString string) bool {
	ServerPublic := util.GetRSA_PKIXPublicKey("serverpub2048.txt")
	//	signature := []byte(data)
	hashed, _ := hex.DecodeString(hashString)
	err := rsa.VerifyPKCS1v15(ServerPublic, crypto.SHA256, hashed[:], data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return false
	}
	//fmt.Println("Valid sign")
	return true

}
func TestSign() {
	genesisBlock := Block{}
	genesisBlock = Block{0, "123456", 0, calculateBlockHash(genesisBlock), "", "", []byte{}}
	block, _ := generateBlock(genesisBlock, 123, "address")
	fmt.Println("is valid ? ", isServerData(block.Data, block.Hash))

}
