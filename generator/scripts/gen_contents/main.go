package main

import (
	"encoding/base64"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/ParteeLabs/gomiger/generator/helper"
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
	migrationTemplateContent, err := os.ReadFile("./generator/mg/migration.mg.go")
	if err != nil {
		fmt.Println("Error reading template file migration.mg.go:", err)
		return
	}
	migratorTemplateContent, err := os.ReadFile("./generator/mg/migrator.mg.go")
	if err != nil {
		fmt.Println("Error reading template file migrator.mg.go:", err)
		return
	}
	cliTemplateContent, err := os.ReadFile("./generator/mg/cli.mg.go")
	if err != nil {
		fmt.Println("Error reading template file cli.mg.go:", err)
		return
	}

	/// Parse the skeleton then add the templates
	fs := token.NewFileSet()
	skeleton, err := parser.ParseFile(fs, "./generator/mg/skeleton.go", nil, parser.AllErrors)
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
	libContentsFile, err := os.Create("./generator/contents.mg.go")
	if err != nil {
		panic(err)
	}
	defer libContentsFile.Close()

	if err := format.Node(libContentsFile, fs, skeleton); err != nil {
		panic(err)
	}
}
