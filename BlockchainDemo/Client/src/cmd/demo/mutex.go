package demo

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type alertCounter int

var message string = "default message"

// NewAlertCounter creates and returns objects of
// the unexported type alertCounter.
func NewAlertCounter(value int) alertCounter {
	fmt.Println("AlertCounter")
	return alertCounter(value)
}
func printAndSleep(name string) {
	fmt.Println(name + " came ")
	mutex.Lock()
	fmt.Println("----------------------------------------")
	fmt.Println(name + " gets message : " + message)
	message = ".." + name + ".. set message"
	time.Sleep(1 * time.Second)
	fmt.Println("Wake up  " + name)
	fmt.Println("----------------------------------------")
	mutex.Unlock()
}
func TestDemo() int {
	fmt.Println("testDemo")
	c := make(chan string)
	c1 := make(chan string)
	c2 := make(chan string)
	go func(cc chan string) {
		for {
			input := <-cc
			printAndSleep("one" + input)
			fmt.Println("1. Send : [" + input + "] ->")
			c1 <- input
			fmt.Println("2. Send : [" + input + "] ->")
			c1 <- input
		}
		//c1 <- input
	}(c)

	go func(name string) {
		for i := 0; i < 2; i++ {
			printAndSleep(name)

		}
		c2 <- "two"
	}("two")
	go func() {
		for {
			input := <-c1
			fmt.Println("-> 1. Received input : [" + input + "]")
		}

	}()
	go func() {
		for {
			input := <-c1
			fmt.Println("-> 2. Received input : [" + input + "]")
		}

	}()
	count := 0
	for {
		c <- " input a message " + strconv.Itoa(count)
		mess := ""
		fmt.Scanf("%s\n", &mess)
		if mess != "" {
			fmt.Println("break ")
			break
		}
		count++
		//	fmt.Println(mess)
	}
	/*
		for i := 0; i < 2; i++ {
			select {
			case msg1 := <-c1:
				fmt.Println("received:", msg1)
			case msg2 := <-c2:
				fmt.Println("received:", msg2)

			}
		}
	*/
	return 0
}
