//nolint:revive
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

// UpdateStringValue updates the string value
func UpdateStringValue(node *ast.File, targetValue, newValue string) *ast.File {
	fmtedTargetValue := fmt.Sprintf("\"%s\"", targetValue)
	ast.Inspect(node, func(n ast.Node) bool {
		if f, ok := n.(*ast.BasicLit); ok && f.Value == fmtedTargetValue {
			f.Value = fmt.Sprintf("\"%s\"", newValue)
		}
		return true
	})
	return node
}

// ExportFile exports the node to a file
func ExportFile(node *ast.File, fs *token.FileSet, path string) error {
	//nolint:gosec
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}()
	if err := format.Node(file, fs, node); err != nil {
		return fmt.Errorf("failed to format ast.File node to file: %w", err)
	}
	return nil
}
