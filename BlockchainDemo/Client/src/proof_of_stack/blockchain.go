package proof_of_stack

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"myutil"
	"os"
	"strconv"
	"sync"
)

// Blockchain is a series of validated Blocks
var Blockchain []Block
var tempBlocks []Block

var ServerPublic *rsa.PublicKey

// candidateBlocks handles incoming blocks for validation
var candidateBlocks = make(chan Block)

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)

var mutex = &sync.Mutex{}

// validators keeps track of open validators and balances
var validators = make(map[string]int)

var Version float64

func Update(datafile string) {
	blocks, _ := json.MarshalIndent(Blockchain, "", "  ")
	myutil.WriteByte(datafile, blocks)
}
func Add(datafile string, block Block) {
	blockJson, _ := json.Marshal(block) //json.MarshalIndent(block, "", "  ")
	row := strconv.Itoa(block.Index) + "=" + string(blockJson)
	//fmt.Println("row: "+row+"->", block.Index)
	myutil.WriteAppendFile(datafile, row)
}
func AddAll(datafile string, blocks []Block) {
	myutil.DeleteFile(datafile)
	for _, block := range blocks {
		Add(datafile, block)
	}

}

func JsonToBlockIf(jsondata []byte, val interface{}) {
	//var result reflect.ValueOf(val)
	json.Unmarshal(jsondata, val)
	//fmt.Println(result)
}
func ReadJson(filePath string) []Block {
	var dao *BlockFile = NewInstance()
	myutil.LoadFileToJson(filePath, dao)
	return dao.Data
}

func Validate(inputJson string) bool {

	val, _ := IsBlockChainValid(inputJson)
	return val
}
func prints(input ...interface{}) {
	mutex.Lock()
	fmt.Println(input)
	mutex.Unlock()
}
func JsonToBlock(jsondata []byte) []Block {
	var result []Block
	json.Unmarshal(jsondata, &result)
	//fmt.Println(result)
	return result
}

func IsBlockChainValid(inputJson string) (bool, []Block) {
	Blocks := JsonToBlock([]byte(inputJson))
	length := len(Blocks)
	for i := length - 1; i > 0; i-- {
		if !isBlockValid(Blocks[i], Blocks[i-1]) ||
			!isServerData(Blocks[i].Data, Blocks[i].Hash) {
			fmt.Println("Block chain is invalid !!!")
			return false, Blocks
		}
	}
	return true, Blocks
}
func isServerData(data []byte, hashString string) bool {

	hashed, _ := hex.DecodeString(hashString)

	err := rsa.VerifyPKCS1v15(ServerPublic, crypto.SHA256, hashed[:], data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err, ServerPublic)
		return false
	}
	return true

}
func LoadServerPublicKey() {
	ServerPublic = myutil.GetRSA_PKIXPublicKey("src/keys/serverpub2048.txt")
}

func LoadVersion() {
	Version = myutil.LoadVersion("verison.txt")
}
func UpdateVersion(version float64) {
	//version += 0.001
	v := fmt.Sprintf("%.3f", version)
	myutil.WriteByte("version.txt", []byte(v))
}
