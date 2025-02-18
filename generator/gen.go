package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

var (
	fs = token.NewFileSet()
)

func LoadTemplates() (migrationNode *ast.File, migratorNode *ast.File, err error) {
	migrationNode, err = parser.ParseFile(fs, "", MigrationScriptTemplate, parser.AllErrors)
	if err != nil {
		return
	}
	migratorNode, err = parser.ParseFile(fs, "", MigratorTemplate, parser.AllErrors)
	if err != nil {
		return
	}
	return
}
