package osexitchecker

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer missing godoc.
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "checks if os.Exit() is using in main() function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			if funcDecl, ok := node.(*ast.FuncDecl); ok {
				if funcDecl.Name.String() != "main" {
					return true
				}

				ast.Inspect(funcDecl, func(node ast.Node) bool {
					if exprStmt, ok := node.(*ast.ExprStmt); ok {
						if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
							if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								x, ok := selExpr.X.(*ast.Ident)
								if !ok {
									return true
								}
								if x.Name == "os" && selExpr.Sel.Name == "Exit" {
									pass.Reportf(selExpr.Pos(), "os.Exit called in main.main()")
								}
							}

						}
					}

					return true
				})
			}

			return true
		})
	}

	return nil, nil
}
