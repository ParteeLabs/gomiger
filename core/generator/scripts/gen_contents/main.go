//go:build ignore

// Package main provides a generator script that creates the contents.mg.go file.
// This file is responsible for embedding template files into the generator package.
//
// The script reads template files (migration, migrator, and CLI templates),
// updates a skeleton file with the template contents encoded in base64,
// and outputs the result as contents.mg.go.
//
// The generator modifies the package name of the skeleton file to "generator"
// and replaces placeholder string literals with the actual template contents.
//
// Template files processed:
//   - migration.mg.go: Template for migration scripts
//   - migrator.mg.go: Template for the migrator implementation
//   - cli.mg.go: Template for CLI interface
//
// The output file contents.mg.go is formatted according to Go standards
// before being written to disk.
package main

import (
	"encoding/base64"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/ParteeLabs/gomiger/core/generator/helper"
)

// main is the script of the generator for `contents.mg.go`.
// The generation process is as follows:
//  1. Read the template files.
//  2. Parse the skeleton file ./generator/mg/skeleton.go.
//  3. Update the package name of the skeleton to "generator".
//  4. Update the string literals in the skeleton to the contents of the
//     template files.
//  5. Write the updated skeleton to the output file.
func main() {
	/// Load template contents
	migrationTemplateContent, err := os.ReadFile("./core/generator/mg/migration.mg.go")
	if err != nil {
		fmt.Println("Error reading template file migration.mg.go:", err)
		return
	}
	migratorTemplateContent, err := os.ReadFile("./core/generator/mg/migrator.mg.go")
	if err != nil {
		fmt.Println("Error reading template file migrator.mg.go:", err)
		return
	}
	cliTemplateContent, err := os.ReadFile("./core/generator/mg/cli.mg.go")
	if err != nil {
		fmt.Println("Error reading template file cli.mg.go:", err)
		return
	}

	/// Parse the skeleton then add the templates
	fs := token.NewFileSet()
	skeleton, err := parser.ParseFile(fs, "./core/generator/mg/skeleton.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing template file migrator.mg.go:", err)
		return
	}

	helper.UpdatePackageName(skeleton, "generator")

	ast.Inspect(skeleton, func(n ast.Node) bool {
		if nf, ok := n.(*ast.BasicLit); ok && nf.Value == "`__MIGRATION_SCRIPT_TEMPLATE__`" {
			nf.Value = fmt.Sprintf("`%s`", base64.StdEncoding.EncodeToString(migrationTemplateContent))
		}
		if nf, ok := n.(*ast.BasicLit); ok && nf.Value == "`__MIGRATOR_TEMPLATE__`" {
			nf.Value = fmt.Sprintf("`%s`", base64.StdEncoding.EncodeToString(migratorTemplateContent))
		}
		if nf, ok := n.(*ast.BasicLit); ok && nf.Value == "`__CLI_TEMPLATE__`" {
			nf.Value = fmt.Sprintf("`%s`", base64.StdEncoding.EncodeToString(cliTemplateContent))
		}
		return true
	})

	/// Write the libContentsFile
	libContentsFile, err := os.Create("./core/generator/contents.mg.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := libContentsFile.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}()

	if err := format.Node(libContentsFile, fs, skeleton); err != nil {
		panic(err)
	}
}
