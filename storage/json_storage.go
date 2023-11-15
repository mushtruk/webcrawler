package storage

import (
	"encoding/json"
	"os"
)

type JsonStorage struct {
	Filename    string
	PrettyPrint bool
}

func NewJSONStorage(filename string, pretty bool) *JsonStorage {
	return &JsonStorage{
		Filename:    filename,
		PrettyPrint: pretty,
	}
}

func (s *JsonStorage) Handle(data interface{}) error {
	var jsonData []byte
	var err error

	if s.PrettyPrint {
		jsonData, err = json.MarshalIndent(data, "", " ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(s.Filename, jsonData, 0644)
}
