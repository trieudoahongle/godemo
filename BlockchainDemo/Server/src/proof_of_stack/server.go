package proof_of_stack

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"util"
)

var ValidatorConst = "validator@@"

// pickWinner creates a lottery pool of validators and chooses the validator who gets to forge a block to the blockchain
// by random selecting from the pool, weighted by amount of tokens staked
func pickWinner() {
	time.Sleep(2 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {

		// slightly modified traditional proof of stake algorithm
		// from all validators who submitted a block, weight them by the number of staked tokens
		// in traditional proof of stake, validators can participate without submitting a block to be forged
	OUTER:
		for _, block := range temp {
			// if already in lottery pool, skip
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// lock list of validators to prevent data race
			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// randomly pick winner from lottery pool
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		// add block of winner to blockchain and let all the other nodes know
		for _, block := range temp {
			if block.Validator == lotteryWinner {
				mutex.Lock()
				//	Blockchain = append(Blockchain, block)
				addBlock(block)
				validateResultCount = 0
				mutex.Unlock()
				for _ = range validators {
					//	fmt.Println("Winner : ", lotteryWinner)
					announcements <- "\nwinning validator: " + lotteryWinner + "\n"
				}
				break
			}
		}
	}

	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}
func handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			//fmt.Println("Announce content:" + msg)
			io.WriteString(conn, msg)
		}
	}()
	// validator address
	var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	//o.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	scanBalance.Scan()
	text := scanBalance.Text()
	fmt.Println(text)
	io.WriteString(conn, "Got test message:"+text+"\nEnter token balance:")
	for scanBalance.Scan() {
		text = strings.Replace(scanBalance.Text(), "\n", "", -1)
		fmt.Println(text)
		balance, err := strconv.Atoi(text)
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		//t := time.Now()
		//Get address
		address = "address_" + text //calculateHash(t.String())
		validators[address] = balance
		//	fmt.Println("Validators are: ", validators)
		break
	}

	io.WriteString(conn, "\nEnter a new BPM:\n")

	scanBPM := bufio.NewScanner(conn)

	go func() {
		for {
			// take in BPM from stdin and add it to blockchain after conducting necessary validation
			for scanBPM.Scan() {
				text := scanBPM.Text()
				if strings.Contains(text, ValidatorConst) {
					content := strings.Replace(text, ValidatorConst, "", -1)
					conArr := strings.SplitN(content, "=", 2)
					switch conArr[0] {
					case "validate_result":
						reArr := strings.SplitN(conArr[1], ":", 2)
						isOk, _ := strconv.ParseBool(reArr[0])
						mutex.Lock()
						if isOk {

							validateResultCount++
						}
						if validateResultCount > 1 {
							for _ = range validators {
								announcements <- "accept=" + reArr[0] + "\n"
							}
						}
						mutex.Unlock()

					case "number_of_lock":
						fmt.Println("Last index of " + address + " is " + conArr[1])
					default:
						log.Printf("Invalid content -> delete validator : " + address)
						delete(validators, address)
						conn.Close()
					}

				} else {
					if len(text) > 0 {
						bpm, err := strconv.Atoi(text)
						//fmt.Println("Got BPM:", bpm)
						// if malicious party tries to mutate the chain with a bad input, delete them as a validator and they lose their staked tokens
						if err != nil {
							log.Printf("%v not a number: %v", scanBPM.Text(), err)
							delete(validators, address)
							conn.Close()
						}

						mutex.Lock()
						oldLastIndex := Blockchain[len(Blockchain)-1]
						mutex.Unlock()

						// create newBlock for consideration to be forged
						newBlock, err := generateBlock(oldLastIndex, bpm, address)
						if err != nil {
							log.Println(err)
							continue
						}
						if isBlockValid(newBlock, oldLastIndex) {
							candidateBlocks <- newBlock
						}
						io.WriteString(conn, "\nEnter a new BPM:\n")
					}
				}
			}
			//For CPU balance
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// simulate receiving broadcast
	for {
		//time.Sleep(30*time.Minute)
		time.Sleep(5 * time.Second)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		//	fmt.Println("Update client: ", output)
		io.WriteString(conn, string(output)+"\n")
	}

}
func initialzeData() {
	LoadVersion()
	if util.IsExist(datafile) {
		Blockchain = ReadJson()
		output, _ := json.Marshal(Blockchain)
		TestJsonMarshal(string(output))

	} else {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", "", []byte{}}
		//spew.Dump(genesisBlock)
		//Blockchain = append(Blockchain, genesisBlock)
		addBlock(genesisBlock)
	}
}
func addBlock(b Block) {
	var length = len(Blockchain)
	isValid := true
	if length > 1 {
		isValid = isBlockValid(b, Blockchain[length-1])
	}
	if isValid {
		Blockchain = append(Blockchain, b)
		UpdateBlocksFile()
		UpdateVersion()
	}
}
func StartProofOfStack() {
	fmt.Println("Load config")
	util.SetEnv("proofofstack.env")
	initialzeData()

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			pickWinner()
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)

	}
}
