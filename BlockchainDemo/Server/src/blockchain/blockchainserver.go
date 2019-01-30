package blockchain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"httpMethod"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Message struct {
	BPM int
}

var difficulty = 1

var mutex = &sync.Mutex{}

func StartBlockChainServer() {
	//	err := godotenv.Load()
	//	if err != nil {
	//	log.Fatal(err)
	//}
	setEnv("app.env")
	go GenerateBlock()
	log.Fatal(run())
}

func GenerateBlock() {

	mutex.Lock()

	Initialize()

	mutex.Unlock()

}

func setEnv(filePath string) {
	f, err := os.Open(filePath)

	if err != nil {
		return
	}
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err == nil {
			lineArr := strings.Split(string(line), "=")
			if len(lineArr) >= 2 {
				os.Setenv(lineArr[0], lineArr[1])
				fmt.Println("Set " + string(line))
			}

		} else {
			break
		}
	}

	f.Close()
}
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	diff := os.Getenv("DIFFICULT")
	difficulty, _ = strconv.Atoi(diff)
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := http.NewServeMux()
	h := httpMethod.NewHttpHandler()
	h.SetHandlerMethod(http.MethodGet, handleGetBlockchain)
	h.SetHandlerMethod(http.MethodPost, handleWriteBlock)
	muxRouter.HandleFunc("/", h.Handler)
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile(datafile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleWriteBlock")
	w.Header().Set("Content-Type", "application/json")

	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//ensure atomicity when creating new block
	mutex.Lock()
	fmt.Println("Start ---------------------")
	newBlock := GoMining(3, Blockchain[len(Blockchain)-1], m.BPM) //generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	fmt.Println("End ---------------------", newBlock)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		Update()
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}
