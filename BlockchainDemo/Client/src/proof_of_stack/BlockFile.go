package proof_of_stack

import (
	"encoding/json"
	"fmt"
	"myutil"
	"strings"
)

type BlockFile struct {
	Data []Block
}

func NewInstance() *BlockFile {
	return &BlockFile{}
}
func (f *BlockFile) SaveFile(filePath string) error {
	myutil.DeleteFile(filePath)
	for _, r := range f.Data {
		row := f.FormatRow(r)
		myutil.WriteAppendFile(filePath, row)
	}
	return nil
}

func (f *BlockFile) LoadFile(FilePath string) {
	fmt.Println("Block LoadFile")
}
func (f *BlockFile) FormatRow(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
func (f *BlockFile) SetRowData(line string) {
	lineArr := strings.SplitN(line, "=", 2)
	var block Block
	json.Unmarshal([]byte(lineArr[1]), &block)
	f.Data = append(f.Data, block)
}
func TestFileDAO(dao myutil.FileDAO) {
	dao.LoadFile("Blockdata_2.txt")
}
func CallDao() {
	var blockDao = NewInstance()
	TestFileDAO(blockDao)
}
