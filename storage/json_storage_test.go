package storage_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/mushtruk/webcrawler/storage"
)

func TestJSONStorage_Handle_Success(t *testing.T) {
	// Setup
	filename := "test_output.json"
	js := storage.NewJSONStorage(filename)
	testData := map[string]string{"key": "value"}

	// Execute
	err := js.Handle(testData)
	if err != nil {
		t.Fatalf("Failed to handle data: %v", err)
	}

	// Verify
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var readData map[string]string
	err = json.Unmarshal(data, &readData)
	if err != nil {
		t.Fatalf("Failed to unmarshal data: %v", err)
	}

	if !reflect.DeepEqual(testData, readData) {
		t.Errorf("Expected %v, got %v", testData, readData)
	}

	os.Remove(filename)
}
