package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/ParteeLabs/gomiger/generator/helper"
)

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
			nf.Value = fmt.Sprintf("`%s`", string(migrationTemplateContent))
			// nf.Value = "`" + string(migrationTemplateContent) + "`"
		}
		if nf, ok := n.(*ast.BasicLit); ok && nf.Value == "`__MIGRATOR_TEMPLATE__`" {
			nf.Value = fmt.Sprintf("`%s`", string(migratorTemplateContent))
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
