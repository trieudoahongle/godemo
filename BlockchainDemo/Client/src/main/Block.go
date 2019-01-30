package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	block.SetHash()
	return block
}
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
func MatchDifficult(byteArr []byte, difficult int) bool {
	//checkStr := strings.Repeat("0", difficult)
	//checkByte := []byte(checkStr)
	for i := 0; i < difficult; i++ {
		if byteArr[i] != 0 {
			return false
		}
	}
	return true
}

func (b *Block) Mining() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	//difficult is global variable
	fmt.Println("Start mining difficult input:", difficult)

	//start := []byte(StartStr)
	for {

		str := fmt.Sprintf("%d", b.Nonce)
		nonce := []byte(str)
		headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, nonce, timestamp}, []byte{})
		hash := sha256.Sum256(headers)
		b.Hash = hash[:]
		if MatchDifficult(b.Hash, difficult) {
			break
		}
		b.Nonce++
	}
	fmt.Println("End  mining difficult input:", b.Nonce)
}
