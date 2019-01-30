package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func WriteByte(output string, data []byte) {

	ioutil.WriteFile(output, data, 0644)
}
func IsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func SetEnv(filePath string) {
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
