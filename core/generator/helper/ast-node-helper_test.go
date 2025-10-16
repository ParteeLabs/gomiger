package helper_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ParteeLabs/gomiger/core/generator/helper"
)

func TestUpdatePackageName(t *testing.T) {
	tests := []struct {
		name       string
		oldPkgName string
		newPkgName string
		expected   string
	}{
		{
			name:       "update normal package name",
			oldPkgName: "oldpkg",
			newPkgName: "newpkg",
			expected:   "newpkg",
		},
		{
			name:       "update to empty package name",
			oldPkgName: "oldpkg",
			newPkgName: "",
			expected:   "",
		},
		{
			name:       "update main package",
			oldPkgName: "main",
			newPkgName: "helper",
			expected:   "helper",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ast.File{
				Name: ast.NewIdent(tt.oldPkgName),
			}

			got := helper.UpdatePackageName(node, tt.newPkgName)

			if got.Name.Name != tt.expected {
				t.Errorf("UpdatePackageName() = %v, want %v", got.Name.Name, tt.expected)
			}
		})
	}
}

func TestUpdateFuncName(t *testing.T) {
	tests := []struct {
		name           string
		sourceCode     string
		targetFuncName string
		newFuncName    string
		expectedCount  int // number of functions with the new name after update
	}{
		{
			name: "update single function name",
			sourceCode: `package main
func OldFunc() {}
func AnotherFunc() {}`,
			targetFuncName: "OldFunc",
			newFuncName:    "NewFunc",
			expectedCount:  1,
		},
		{
			name: "update multiple functions with same name",
			sourceCode: `package main
func DuplicateFunc() {}
func DuplicateFunc(x int) {}
func OtherFunc() {}`,
			targetFuncName: "DuplicateFunc",
			newFuncName:    "RenamedFunc",
			expectedCount:  2,
		},
		{
			name: "target function not found",
			sourceCode: `package main
func ExistingFunc() {}`,
			targetFuncName: "NonExistentFunc",
			newFuncName:    "NewFunc",
			expectedCount:  0,
		},
		{
			name: "update method with receiver",
			sourceCode: `package main
type MyStruct struct{}
func (m MyStruct) OldMethod() {}
func OldMethod() {}`,
			targetFuncName: "OldMethod",
			newFuncName:    "NewMethod",
			expectedCount:  2, // both method and function should be renamed
		},
		{
			name:           "empty source code",
			sourceCode:     "package main",
			targetFuncName: "AnyFunc",
			newFuncName:    "NewFunc",
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "", tt.sourceCode, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse source code: %v", err)
			}

			updatedNode := helper.UpdateFuncName(node, tt.targetFuncName, tt.newFuncName)

			// Count functions with the new name
			count := 0
			ast.Inspect(updatedNode, func(n ast.Node) bool {
				if f, ok := n.(*ast.FuncDecl); ok && f.Name.Name == tt.newFuncName {
					count++
				}
				return true
			})

			if count != tt.expectedCount {
				t.Errorf("UpdateFuncName() resulted in %d functions with name '%s', want %d",
					count, tt.newFuncName, tt.expectedCount)
			}

			// Ensure the original node is returned (modified in place)
			if updatedNode != node {
				t.Error("UpdateFuncName() should return the same node instance")
			}
		})
	}
}

func TestUpdateStringValue(t *testing.T) {
	tests := []struct {
		name          string
		sourceCode    string
		targetValue   string
		newValue      string
		expectedCount int // number of strings with the new value after update
	}{
		{
			name: "update single string literal",
			sourceCode: `package main
const Message = "hello"
const Other = "world"`,
			targetValue:   "hello",
			newValue:      "hi",
			expectedCount: 1,
		},
		{
			name: "update multiple identical strings",
			sourceCode: `package main
const A = "duplicate"
const B = "duplicate"
const C = "different"`,
			targetValue:   "duplicate",
			newValue:      "updated",
			expectedCount: 2,
		},
		{
			name: "target string not found",
			sourceCode: `package main
const Message = "hello"`,
			targetValue:   "nonexistent",
			newValue:      "new",
			expectedCount: 0,
		},
		{
			name: "update empty string",
			sourceCode: `package main
const Empty = ""
const NotEmpty = "text"`,
			targetValue:   "",
			newValue:      "filled",
			expectedCount: 1,
		},
		{
			name: "update string with special characters",
			sourceCode: `package main
const Special = "hello world"`,
			targetValue:   "hello world",
			newValue:      "hi universe",
			expectedCount: 1,
		},
		{
			name:          "no string literals in source",
			sourceCode:    "package main\nfunc Test() {}",
			targetValue:   "any",
			newValue:      "new",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "", tt.sourceCode, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse source code: %v", err)
			}

			updatedNode := helper.UpdateStringValue(node, tt.targetValue, tt.newValue)

			// Count strings with the new value
			count := 0
			expectedLiteral := `"` + tt.newValue + `"`
			ast.Inspect(updatedNode, func(n ast.Node) bool {
				if f, ok := n.(*ast.BasicLit); ok && f.Value == expectedLiteral {
					count++
				}
				return true
			})

			if count != tt.expectedCount {
				t.Errorf("UpdateStringValue() resulted in %d strings with value '%s', want %d",
					count, tt.newValue, tt.expectedCount)
			}

			// Ensure the original node is returned (modified in place)
			if updatedNode != node {
				t.Error("UpdateStringValue() should return the same node instance")
			}
		})
	}
}

func TestExportFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "ast-helper-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		sourceCode    string
		filename      string
		expectError   bool
		errorContains string
	}{
		{
			name: "export simple go file",
			sourceCode: `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			filename:    "hello.go",
			expectError: false,
		},
		{
			name:        "export empty package",
			sourceCode:  `package empty`,
			filename:    "empty.go",
			expectError: false,
		},
		{
			name: "export complex structure",
			sourceCode: `package complex

type MyStruct struct {
	Field1 string
	Field2 int
}

func (m MyStruct) Method() string {
	return m.Field1
}

const (
	Const1 = "value1"
	Const2 = 42
)`,
			filename:    "complex.go",
			expectError: false,
		},
		{
			name:          "export to invalid path",
			sourceCode:    `package main`,
			filename:      "/invalid/path/that/does/not/exist/file.go",
			expectError:   true,
			errorContains: "failed to create file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "", tt.sourceCode, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse source code: %v", err)
			}

			var filePath string
			if tt.expectError && strings.Contains(tt.filename, "/invalid/") {
				// Use the invalid path as-is for error testing
				filePath = tt.filename
			} else {
				filePath = filepath.Join(tempDir, tt.filename)
			}

			err = helper.ExportFile(node, fset, filePath)

			if tt.expectError {
				if err == nil {
					t.Error("ExportFile() expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("ExportFile() error = %v, want error containing '%s'", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ExportFile() unexpected error: %v", err)
				return
			}

			// Verify file was created and contains expected content
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Error("ExportFile() did not create the expected file")
				return
			}

			// Read the file back and verify it's valid Go code
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read exported file: %v", err)
				return
			}

			// Parse the exported content to ensure it's valid
			_, err = parser.ParseFile(token.NewFileSet(), "", string(content), parser.ParseComments)
			if err != nil {
				t.Errorf("Exported file contains invalid Go code: %v", err)
			}

			// Verify the package name matches
			if !strings.Contains(string(content), "package") {
				t.Error("Exported file does not contain package declaration")
			}
		})
	}
}
