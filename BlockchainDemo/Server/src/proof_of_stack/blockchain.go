package proof_of_stack

import (
	"encoding/json"
	"fmt"
	//"strings"
	"io/ioutil"
	"myutil"
	"sync"
)

// Blockchain is a series of validated Blocks
var Blockchain []Block
var tempBlocks []Block

// candidateBlocks handles incoming blocks for validation
var candidateBlocks = make(chan Block)

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)

var mutex = &sync.Mutex{}

// validators keeps track of open validators and balances
var validators = make(map[string]int)

var validateResultCount int = 0
var version float64

const datafile = "blockchain.txt"
const versionfile = "version.txt"

func ReadJson() []Block {
	data, err := ioutil.ReadFile(datafile)
	if err != nil {
		return nil
	}
	var result []Block
	json.Unmarshal(data, &result)
	return result
}
func TestJsonMarshal(jsonStr string) {
	var result []Block
	data := []byte(jsonStr)
	fmt.Println(len(jsonStr))
	fmt.Println(len(data))
	json.Unmarshal(data, &result)
	fmt.Println(len(result))
}
func UpdateBlocksFile() {
	blocks, _ := json.MarshalIndent(Blockchain, "", "  ")
	myutil.WriteByte(datafile, blocks)

}
func LoadVersion() {
	version = myutil.LoadVersion(versionfile)
}
func UpdateVersion() {
	version += 0.001
	v := fmt.Sprintf("%.3f", version)
	myutil.WriteByte(versionfile, []byte(v))
}
