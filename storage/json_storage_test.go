package storage_test

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mushtruk/webcrawler/storage"
)

func TestJSONStorage_Handle_Success(t *testing.T) {
	// Setup
	filename := "test_output.json"
	js := storage.NewJSONStorage(filename, false)
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

func TestJSONStorage_PrettyPrintedOutput(t *testing.T) {
	filename := "test_pretty_output.json"
	js := storage.NewJSONStorage(filename, true)

	testData := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
		"key3": []string{"one", "two", "three"},
	}

	err := js.Handle(testData)
	if err != nil {
		t.Fatalf("Failed to handle data: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(data), "\n") {
		t.Error("Expected pretty-printed JSON output with indentation")
	}

	os.Remove(filename)
}

func TestJSONStorage_FileSystemError(t *testing.T) {
	filename := "/non_existent_directory/test_output.json"
	js := storage.NewJSONStorage(filename, false)

	testData := map[string]string{"key": "value"}

	err := js.Handle(testData)
	if err == nil {
		t.Fatal("Expected file system error, got nil")
	}
}
