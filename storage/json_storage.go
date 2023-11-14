package storage

import (
	"encoding/json"
	"os"
)

type JsonStorage struct {
	Filename string
}

func NewJSONStorage(filename string) *JsonStorage {
	return &JsonStorage{
		Filename: filename,
	}
}

func (s *JsonStorage) Handle(data interface{}) error {
	json, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return os.WriteFile(s.Filename, json, 0644)
}
