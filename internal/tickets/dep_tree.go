package tickets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DepTreeNode struct {
	ID       string
	Status   string
	Children []*DepTreeNode
}

func DepTree(dir, partialID string, full bool) (*DepTreeNode, error) {
	ticketsDir := filepath.Join(dir, ".tickets")
	archiveDir := filepath.Join(ticketsDir, "archive")

	path, err := resolveID(partialID, ticketsDir, archiveDir)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fm, _, err := parseFrontmatter(raw)
	if err != nil {
		return nil, err
	}

	id, _ := fm["id"].(string)
	status, _ := fm["status"].(string)

	node := &DepTreeNode{ID: id, Status: status}

	deps := stringSlice(fm["dependencies"])
	for _, depID := range deps {
		if full {
			child, err := DepTree(dir, depID, true)
			if err != nil {
				child = &DepTreeNode{ID: depID, Status: "unknown"}
			}
			node.Children = append(node.Children, child)
		} else {
			depStatus := statusOf(depID, ticketsDir, archiveDir)
			node.Children = append(node.Children, &DepTreeNode{ID: depID, Status: depStatus})
		}
	}

	return node, nil
}

func RenderDepTree(root *DepTreeNode) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s (%s)\n", root.ID, root.Status)
	for i, child := range root.Children {
		renderNode(&b, child, "", i == len(root.Children)-1)
	}
	return b.String()
}

func renderNode(b *strings.Builder, node *DepTreeNode, prefix string, isLast bool) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	}
	fmt.Fprintf(b, "%s%s%s (%s)\n", prefix, connector, node.ID, node.Status)
	for i, child := range node.Children {
		renderNode(b, child, childPrefix, i == len(node.Children)-1)
	}
}
