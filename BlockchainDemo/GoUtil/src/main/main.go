package main

import (
	"fmt"
	"os"
	"time"
)

var p = fmt.Println

func main() {
	fmt.Println("GoUtil")
	argsWithProg := os.Args
	p(argsWithProg)
	argsWithoutProg := os.Args[1:]

	for _, arg := range argsWithoutProg {
		p(arg + "  ")
		time.Sleep(30 * time.Second)
	}

}
