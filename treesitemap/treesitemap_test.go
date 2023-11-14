package treesitemap

import (
	"strings"
	"testing"
)

func TestTreeNodeCreationAndHierarchy(t *testing.T) {
	rootNode := NewTreeNode("http://root-node.test.com", "Root Content")

	childNode := NewTreeNode("http://child-node.test.com", "Child Content")

	rootNode.AddChild(childNode)

	if len(rootNode.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(rootNode.Children))
	}

	if rootNode.Children[0].URL != "http://child-node.test.com" {
		t.Errorf("Expected child URL to be 'http://child.com', got %s", rootNode.Children[0].URL)
	}
}

func TestTreeSerialization(t *testing.T) {
	rootNode := NewTreeNode("http://root.com", "Root Content")
	childNode := NewTreeNode("http://child.com", "Child Content")
	rootNode.AddChild(childNode)

	jsonData, err := SerializeTree(rootNode)

	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	if !strings.Contains(jsonData, "http://root.com") || !strings.Contains(jsonData, "http://child.com") {
		t.Errorf("Serialized data does not contain expected URLs")
	}
}
