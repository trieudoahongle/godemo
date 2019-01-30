package main

import (
	//"cars"
	//"cmd/demo"
	"fmt"
	//"request"
	"log"
	"net"
	"proof_of_stack"
	//"singleton"
	"strings"
	"time"
)

var difficult int = 2

func printPtr(i *int) {
	fmt.Println(i)
	fmt.Println(&i)
	fmt.Println(*i)
}

func swap(i *int, j *int) {
	t := *i
	*i = *j
	*j = t
}
func swap2(i int, j int) {
	swap(&i, &j)
}

func testRoofOfStack() {
	var bpmArr []string
	bpmArr = append(bpmArr, "3")
	bpmArr = append(bpmArr, "5")

	var tokenArr []string
	tokenArr = append(tokenArr, "2")
	tokenArr = append(tokenArr, "4")

	mapConn := make(map[string]net.Conn)
	mapChanel := make(map[string]chan []string)
	proof_of_stack.LoadServerPublicKey()

	for _, T := range tokenArr {
		conn, err := net.Dial("tcp", "127.0.0.1:9000")
		if err != nil {
			log.Println(err)
			return
		}

		mapConn[T] = conn
		chanel := make(chan []string)
		mapChanel[T] = chanel
		proof_of_stack.CallTCLForForging(T, mapConn[T])
	}

	for i := 0; i < len(tokenArr); i++ {
		T := tokenArr[i]
		go proof_of_stack.SendBPMChanel(mapChanel[T], mapConn[T])
		go proof_of_stack.ReadResult(T, mapConn[T])
	}

	for {
		for _, T := range tokenArr {
			mapChanel[T] <- bpmArr
			//request.SendBPM(bpmArr, mapConn[T])
		}
		fmt.Println("\nSent all. -> waiting result. -> sleep ")
		time.Sleep(3 * time.Second)
		fmt.Println("\nWake up and Send ")

	}

	fmt.Scanln()
	for _, T := range tokenArr {
		defer mapConn[T].Close()
	}
}
func main() {
	fmt.Println(time.Now().UnixNano(), " start hello world hello main")

	/*
		car := cars.NewCar()
		car.HonkTheHorn()
		fmt.Println(time.Now().UnixNano(), " end")
		r := singleton.Repository()
		fmt.Println(r)
		r.Set("key", "value to set 1")
		item, _ := r.Get("key")
		fmt.Println("item 1 = ", item)
	*/
	//testBlock()

	//	demo.CheckHash("000000000000")

	/*
		r2 := singleton.Repository()
		fmt.Println(r2)
		r2.Set("key", "value to set 2")
		item2, _ := r2.Get("key")
		fmt.Println("item 2 = ", item2)
		client := request.NewBasicAuthClient("user", "pwd")
		client.GetTodo(1)

		b := demo.NewBankAdapter()
		demo.PrintAdapter(b)
		p := demo.NewPaypalAdapter()
		demo.PrintAdapter(p)
		request.CallTCL()
	*/
	input := "one two one aa"

	// Get last index, searching from right to left.
	result := strings.SplitN(input, " ", 2)
	fmt.Println("Index: ", strings.Index(input, " "))
	fmt.Println("LastIndex: ", strings.LastIndex(input, " "))
	for _, str := range result {
		fmt.Println(str)
	}
	//proof_of_stack.ReadJson("Blockdata_2.txt")
	//fmt.Println("Get file: ", blocks)

	testRoofOfStack()
	fmt.Println("done")
	/*
		var i int = 6
		var j int = 1
		printPtr(&i)
		fmt.Println("i=", i, "j=", j)
		swap(&i, &j)
		fmt.Println("i=", i, "j=", j)
		swap2(i, j)
		fmt.Println("i=", i, "j=", j)*/
	//demo.TestDemo()

}
