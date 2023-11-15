package treesitemap

type TreeNode struct {
	URL      string
	Content  interface{}
	Children []TreeNode
}

func NewTreeNode(url, content string) *TreeNode {
	return &TreeNode{
		URL:      url,
		Content:  content,
		Children: make([]TreeNode, 0),
	}
}

func (tn *TreeNode) AddChild(child *TreeNode) {
	tn.Children = append(tn.Children, *child)
}
