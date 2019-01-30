package myutil

import (
	"encoding/json"
	"fmt"
)

// Checkouter checkouts order
type FileDAO interface {
	// Pay from email to email this amount
	SaveFile(FilePath string) error
	FormatRow(data interface{}) string
	LoadFile(FilePath string)
	SetRowData(json string)
}

type JsonFile struct {
	Data []interface{}
}

func NewInstance() *JsonFile {
	return &JsonFile{}
}
func (f *JsonFile) SaveFile(filePath string) error {
	DeleteFile(filePath)
	for _, r := range f.Data {
		row := f.FormatRow(r)
		WriteAppendFile(filePath, row)
	}
	return nil
}

func (f *JsonFile) LoadFile(FilePath string) {
	fmt.Println("LoadFile")
}
func (f *JsonFile) FormatRow(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
func (f *JsonFile) SetRowData(json string) {
	fmt.Println("SetRowData " + json)
}
func TestFileDAO(dao FileDAO) {
	dao.LoadFile("Blockdata_2.txt")
}
func CallDao() {
	var jsonDao = NewInstance()
	TestFileDAO(jsonDao)
}
