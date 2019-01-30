package blockchain

import (
	//"crypto/sha256"
	//"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"util"
)

var c1 = make(chan string)
var c2 = make(chan string)
var blockmutex = &sync.Mutex{}

const datafile = "blockchain.json"

var Blockchain []Block
var FoundBlock chan Block

var FoundNonce bool

func isFoundNounce() bool {
	blockmutex.Lock()
	defer blockmutex.Unlock()
	return FoundNonce
}

func setFoundBlock(block Block) bool {
	blockmutex.Lock()
	defer blockmutex.Unlock()
	if FoundNonce {
		fmt.Println(" Noucce is found !")
		return false
	}
	FoundNonce = true
	FoundBlock <- block
	return true
}
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func ReadJson(input string) []Block {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return nil
	}
	var result []Block
	json.Unmarshal(data, &result)
	fmt.Println(result)
	return result
}

func ReadFile(input string) []byte {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return nil
	}
	return data
}

func Initialize() {
	//if !util.IsExist(datafile) {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, "", t.String(), 0, difficulty, 0, calculateHash(genesisBlock)}

	Blockchain = append(Blockchain, genesisBlock)
	newlock := genesisBlock
	for i := 1; i < 5; i++ {
		newlock = generateBlock(newlock, i)
		Blockchain = append(Blockchain, newlock)
	}
	Update()
	//} else {
	//	Blockchain = ReadJson(datafile)
	//}
}
func Update() {
	blocks, _ := json.MarshalIndent(Blockchain, "", "  ")
	util.WriteByte(datafile, blocks)

}

func GoMining(numberMiner int, oldBlock Block, BPM int) Block {
	FoundNonce = false
	FoundBlock = make(chan Block)
	for i := 0; i < numberMiner; i++ {
		go Mining("minner "+strconv.Itoa(i), oldBlock, BPM)
	}

	fmt.Println(" End of mining nonce ", FoundBlock)
	for i := 0; i < numberMiner; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1, FoundBlock)
			//return getWinnerBlock(msg1, oldBlock, BPM)
		case block := <-FoundBlock:
			fmt.Println("received", block, FoundBlock)
			return block
		}
	}
	return <-FoundBlock
}

func getWinnerBlock(result string, oldBlock Block, BPM int) Block {
	arrStr := strings.Split(result, ",")
	fmt.Println(" Winner  ", arrStr[0])
	nonce, _ := strconv.Atoi(arrStr[1])
	fmt.Println(" Nonce  ", nonce)
	return generateBlock(oldBlock, BPM)
}

func getBlock(result Block) Block {
	return result
}

func Mining(name string, oldBlock Block, BPM int) int {
	/*
		sleepMilis := rand.Intn(3)
		//time.Sleep(sleepMilis * time.)
		fmt.Println(name+" start sleep ", sleepMilis)
		for i := 0; i < sleepMilis; i++ {
			time.Sleep(1 * time.Millisecond)
		}
		fmt.Println(name + " start wake up ")
	*/
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	fmt.Println(name + " starts mining ")
	for i := 0; ; i++ {
		//	hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = i
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			if FoundNonce {
				fmt.Println(name + " stops mining ")
				c2 <- name
				break
			}
			continue
		} else {

			fmt.Println(calculateHash(newBlock), " work done!", name)
			newBlock.Hash = calculateHash(newBlock)
			if !setFoundBlock(newBlock) {
				fmt.Println(name + " stops mining")
			} else {
				fmt.Println(name+" found nonce ", FoundBlock)
				c1 <- name + "," + strconv.Itoa(i)
			}
			break
		}

	}
	return 1
}
