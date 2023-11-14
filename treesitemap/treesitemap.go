package treesitemap

import (
	"encoding/json"
	"os"
)

func SerializeTree(node *TreeNode) (string, error) {

	jsonData, err := json.Marshal(node)

	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func WriteToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}
