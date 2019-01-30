package request

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"myutil"
	"net"
	"net/http"
	"strings"
)

const baseURL string = "http://localhost:8080"

type Client struct {
	Username string
	Password string
}

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
	fmt.Println("call ", url)
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
	fmt.Println("doRequest : ")
	req.SetBasicAuth(s.Username, s.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Reponse err : ", err)
	fmt.Println("Reponse body : ", body)
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

func ScanResult(balance string, conn net.Conn) {
	fmt.Println("Scan "+balance, conn)

	scanResult := bufio.NewScanner(conn)

	for {
		// take in BPM from stdin and add it to blockchain after conducting necessary validation
		for scanResult.Scan() {
			text := scanResult.Text()

			if strings.Contains(text, "[") {
				filename := "Blockdata_" + balance + ".txt"
				fmt.Println("Write file: " + filename + "\nContent:" + text)
				util.WriteAppendFile(filename, text)
			}
			if strings.Contains(text, "validator") {
				fmt.Println(" " + balance + " -> " + text)
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
	println("Send content:" + cmd)
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
	println(string(buf[:n]))
	return nil
}
