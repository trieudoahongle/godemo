package proof_of_stack

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	//"util"
	"myutil"
)

const baseURL string = "http://localhost:8080"

type Client struct {
	Username string
	Password string
}

var ValidatorConst = "validator@@"

func NewBasicAuthClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
	}
}

type Todo struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func (s *Client) GetTodo(id int) (*Todo, error) {
	url := fmt.Sprintf(baseURL+"/view/%s/%d", s.Username, id)
	prints("call ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data Todo
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
func (s *Client) doRequest(req *http.Request) ([]byte, error) {
	prints("doRequest : ")
	req.SetBasicAuth(s.Username, s.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	prints("Reponse err : ", err)
	prints("Reponse body : ", body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func CallTCL() {
	log.SetFlags(log.Lshortfile)

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	err = sendCommand("ehlo", conn)

	if err != nil {
		log.Println(err)
		return
	}
	err = sendCommand("mail from:", conn)

	if err != nil {
		log.Println(err)
		return
	}
}

//block chain input a string of int value : balance & scanBPM
//In order to not convert -> use string
func CallTCLForForging(balance string, conn net.Conn) {
	log.SetFlags(log.Lshortfile)

	err := sendCommand("Hello this is test message of "+balance, conn)

	if err != nil {
		log.Println(err)
		return
	}

	err = sendCommand(balance, conn)

	if err != nil {
		log.Println(err)
		return
	}

}

func SendBPM(scanBPM []string, conn net.Conn) {

	for _, bpm := range scanBPM {
		err := sendTCLCommand(bpm, conn)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func SendBPMChanel(getChanelBPM <-chan []string, conn net.Conn) {
	go func() {
		for {
			scanBPM := <-getChanelBPM
			//fmt.Println("scanBPM len:", len(scanBPM))
			for _, bpm := range scanBPM {
				err := sendTCLCommand(bpm, conn)

				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()
}
func rotateFileName(filePath string) {
	MAX_BACKUP := 1
	bkFileName := filePath + strconv.Itoa(MAX_BACKUP)
	os.Remove(bkFileName)
	bkCurr := filePath + "_bk1"
	for i := MAX_BACKUP; i > 1; i-- {
		bkCurr = filePath + "_bk" + strconv.Itoa(i-1)
		if myutil.IsExist(bkCurr) {
			bkOld := filePath + "_bk" + strconv.Itoa(i)
			os.Rename(bkCurr, bkOld)
		}
	}
	os.Rename(filePath, bkCurr)
}
func ReadResult(balance string, conn net.Conn) {
	r := bufio.NewReader(conn)
	prints("Scan "+balance, conn)

	filename := "Blockdata_" + balance + ".txt"
	var lastIndex int
	//Version := myutil.LoadVersion("verison.txt")
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		msg = strings.Replace(msg, "\n", "", -1)
		TreatMessage(balance, msg, filename, conn, &lastIndex)
	}
}
func TreatMessage(balance string, text string, filename string, conn net.Conn, lastIndex *int) {
	isValid := true

	if strings.Contains(text, "[") {

		prints("Write file: " + filename)

		//blocks := JsonToBlock([]byte(text))
		var blocks []Block
		//JsonToBlockIf([]byte(text), &blocks)
		fmt.Println("Length of receive text :", len(text))
		//util.WriteByte("receiveData.txt", []byte(text))
		isValid, blocks = IsBlockChainValid(text)
		if isValid {
			length := len(blocks)
			//	fmt.Println(" -> length :", length)
			*lastIndex = blocks[length-1].Index
			if myutil.IsExist(filename) {
				rotateFileName(filename)
			}
			AddAll(filename, blocks)

		} //	util.WriteAppendFile(filename, text)
	}
	if strings.Contains(text, "validator") {
		prints(balance + " -> Validate :" + text)
		//	isValid := Validate()
		returnStr := strconv.FormatBool(isValid)
		err := sendTCLCommand(ValidatorConst+"validate_result="+returnStr+":"+text, conn)

		if err != nil {
			log.Println(err)
		}
	}

	if strings.Contains(text, "accept") {
		prints(balance + " -> Save blockchain <------------")

		//readBlocks := ReadJson(filename)
		//length := len(readBlocks)
		if *lastIndex >= 0 {
			err := sendTCLCommand(ValidatorConst+"number_of_lock="+strconv.Itoa(*lastIndex), conn)

			if err != nil {
				log.Println(err)
			}
		}
	}
}
func ScanResult(balance string, conn net.Conn) {
	prints("Scan "+balance, conn)

	scanResult := bufio.NewScanner(conn)
	filename := "Blockdata_" + balance + ".txt"
	//	tmpFilename := filename + ".tmp"
	os.Remove(filename)
	//	os.Remove(tmpFilename)

	for {
		// take in BPM from stdin and add it to blockchain after conducting necessary validation
		isValid := true
		var lastIndex int

		for scanResult.Scan() {
			text := scanResult.Text()

			if strings.Contains(text, "[") {

				prints("Write file: " + filename)

				//blocks := JsonToBlock([]byte(text))
				var blocks []Block
				//JsonToBlockIf([]byte(text), &blocks)
				fmt.Println("Length of receive text :", len(text))
				myutil.WriteByte("receiveData.txt", []byte(text))
				isValid, blocks = IsBlockChainValid(text)
				if isValid {
					fmt.Println("Block is :", blocks)
					length := len(blocks)
					fmt.Println(" -> length :", length)
					lastIndex = blocks[length-1].Index
					if myutil.IsExist(filename) {
						rotateFileName(filename)
					}
					AddAll(filename, blocks)

				} //	util.WriteAppendFile(filename, text)
			}
			if strings.Contains(text, "validator") {
				prints(balance + " -> Validate :" + text)
				//	isValid := Validate()
				returnStr := strconv.FormatBool(isValid)
				err := sendTCLCommand(ValidatorConst+"validate_result="+returnStr+":"+text, conn)

				if err != nil {
					log.Println(err)
				}
			}

			if strings.Contains(text, "accept") {
				prints(balance + " -> Save blockchain <------------")

				//readBlocks := ReadJson(filename)
				//length := len(readBlocks)
				if lastIndex >= 0 {
					err := sendTCLCommand(ValidatorConst+"number_of_lock="+strconv.Itoa(lastIndex), conn)

					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
func sendTCLCommand(cmd string, conn net.Conn) error {
	println("Send TCP content:" + cmd)
	n, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Println(n, err)
		return err
	}
	return nil
}

func sendCommand(cmd string, conn net.Conn) error {
	prints("Send content:" + cmd)
	n, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Println(n, err)
		return err
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return err
	}
	prints(string(buf[:n]))
	return nil
}
