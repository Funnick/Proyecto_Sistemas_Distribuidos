package chord

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type DataBasePlatform interface {
	Get([]byte) (string, bool)
	//GetByFun(string) (string, bool)
	GetAll() ([]string, bool)
	Set([]byte, string) bool
	Update([]byte, string) bool
	Delete([]byte) bool
}

// Base de Datos para la plataforma
type DataBasePl struct {
	fileName string
}

// Cada informacion es un par key-data
type RowData struct {
	Key  []byte
	Data string
}

func NewDataBase(fileName string) *DataBasePl {
	db := &DataBasePl{
		fileName: fileName,
	}

	file, err := os.OpenFile(db.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer file.Close()

	return db
}

func (db *DataBasePl) readAll() ([]RowData, error) {
	file, err := os.Open(db.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return []RowData{}, nil
	}

	var rows []RowData
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&rows); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return rows, nil
}

func (db *DataBasePl) writeAll(rows []RowData) error {
	file, err := os.Create(db.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(rows); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (db *DataBasePl) Get(key []byte) (string, bool) {
	rows, err := db.readAll()
	if err != nil {
		return "", false
	}

	for _, elem := range rows {
		if bytes.Equal(elem.Key, key) {
			return elem.Data, true
		}
	}

	return "", false
}

func (db *DataBasePl) GetByFun(fun string) ([]string, bool) {
	rows, err := db.readAll()
	if err != nil {
		return []string{}, false
	}

	data := make([]string, 0)
	for _, elem := range rows {
		if SearchString(elem.Data, fun) != -1 {
			data = append(data, elem.Data)
		}
	}
	if len(data) > 0 {
		return data, true
	}
	return data, false
}

func (db *DataBasePl) GetAll() ([]string, bool) {
	rows, err := db.readAll()
	if err != nil {
		return nil, false
	}

	var data []string = make([]string, len(rows))

	for i, elem := range rows {
		data[i] = elem.Data
	}

	return data, true
}

func (db *DataBasePl) Set(vKey []byte, vData string) bool {
	rows, err := db.readAll()
	if err != nil {
		return false
	}

	for _, elem := range rows {
		if bytes.Equal(elem.Key, vKey) {
			return false
		}
	}

	var newRows []RowData = append(rows, RowData{Key: vKey, Data: vData})
	err = db.writeAll(newRows)
	if err != nil {
		return false
	}

	return true
}

func (db *DataBasePl) Update(vKey []byte, vData string) bool {
	rows, err := db.readAll()
	if err != nil {
		return false
	}

	for i, elem := range rows {
		if bytes.Equal(elem.Key, vKey) {

			newRows := append(rows[:i], rows[i+1:]...)
			newRows = append(newRows, RowData{Key: vKey, Data: vData})
			err = db.writeAll(newRows)
			if err != nil {
				return false
			}

			return true
		}
	}

	return false
}

func (db *DataBasePl) Delete(vKey []byte) bool {
	rows, err := db.readAll()
	if err != nil {
		return false
	}

	for i, elem := range rows {
		if bytes.Equal(elem.Key, vKey) {
			newRows := append(rows[:i], rows[i+1:]...)
			err = db.writeAll(newRows)
			if err != nil {
				return false
			}

			return true
		}
	}

	return false
}
