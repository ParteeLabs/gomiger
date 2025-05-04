package helper

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

// UpdatePackageName updates the package name
func UpdatePackageName(node *ast.File, packageName string) *ast.File {
	node.Name.Name = packageName
	return node
}

// UpdateFuncName updates the function name
func UpdateFuncName(node *ast.File, targetFuncName, newFuncName string) *ast.File {
	ast.Inspect(node, func(n ast.Node) bool {
		if f, ok := n.(*ast.FuncDecl); ok && f.Name.Name == targetFuncName {
			f.Name.Name = newFuncName
		}
		return true
	})
	return node
}

// ExportFile exports the node to a file
func ExportFile(node *ast.File, fs *token.FileSet, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	if err := format.Node(file, fs, node); err != nil {
		return fmt.Errorf("failed to format ast.File node to file: %w", err)
	}
	return nil
}
