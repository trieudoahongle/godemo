package myutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func WriteByte(output string, data []byte) {

	ioutil.WriteFile(output, data, 0644)
}
func WriteAppendFile(output string, data string) {
	if IsExist(output) {
		fileHandle, _ := os.OpenFile(output, os.O_APPEND, 0666)
		writer := bufio.NewWriter(fileHandle)
		defer fileHandle.Close()

		fmt.Fprintln(writer, data)
		writer.Flush()
	} else {
		WriteByte(output, []byte(data+"\n"))
	}
}
func IsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
func DeleteFile(path string) {
	// delete file
	if IsExist(path) {
		var err = os.Remove(path)
		if isError(err) {
			fmt.Println(err)
		}
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
			lineArr := strings.SplitN(string(line), "=", 2)
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
func LoadVersion(filename string) float64 {
	f, err := os.Open(filename)

	if err != nil {
		return 0
	}
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err == nil {
			version, err := strconv.ParseFloat(string(line), 64)
			if err == nil {
				return version
			}
		} else {
			break
		}
	}

	f.Close()

	return 0

}
func LoadFileToJson(filePath string, dao FileDAO) error {
	f, err := os.Open(filePath)
	if err == nil {
		r := bufio.NewReader(f)

		for {
			line, _, err := r.ReadLine()
			//	lineArr := strings.SplitN(string(line), "=", 2)
			if err == nil {
				dao.SetRowData(string(line))
			} else {
				break
			}
		}

	}
	f.Close()
	return nil
}

func LoadFileToJsonIF(filePath string, i interface{}, arrI interface{}) interface{} {
	f, err := os.Open(filePath)
	t := reflect.TypeOf(i)
	var arr []interface{}
	sliceI := reflect.MakeSlice(reflect.TypeOf(arrI), 0, 0).Interface()
	fmt.Println(reflect.TypeOf(sliceI))

	if err == nil {
		r := bufio.NewReader(f)

		for {
			line, _, err := r.ReadLine()

			if err == nil {
				lineArr := strings.SplitN(string(line), "=", 2)
				//dao.SetRowData(string(line))
				//fmt.Println(lineArr[0])
				block := reflect.New(t).Interface()
				json.Unmarshal([]byte(lineArr[1]), &block)
				arr = append(arr, block)
				//sliceI = reflect.Append(sliceI, block)
				sliceI = block
				//fmt.Println(block)
			} else {
				break
			}
		}

	}
	f.Close()
	return arr
}
func fill(i interface{}, arrIf []interface{}) error {
	newLen := len(arrIf)
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %v", v.Type())
	}
	// get the value that the pointer v points to.
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("can't fill non-slice value")
	}
	v.Set(reflect.MakeSlice(v.Type(), newLen, newLen))

	for _, w := range arrIf {
		rv := reflect.ValueOf(w)
		fmt.Println(rv)
		v = reflect.Append(v, rv)
	}
	fmt.Println(v)
	fmt.Println(i)
	return nil
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
