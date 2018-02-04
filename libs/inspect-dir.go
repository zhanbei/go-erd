package libs

import (
	"go/token"
	"os"
	"go/parser"
	"log"
	"fmt"
	"go/ast"
)

func InspectDir(path string) map[string]map[string]NamedType {
	var (
		fset        = token.NewFileSet()
		filter      = func(n os.FileInfo) bool { return true }
		pkgmap, err = parser.ParseDir(fset, path, filter, 0)

		types = make(map[string]map[string]NamedType)
	)

	if err != nil {
		log.Fatal("parser error:", err)
	}

	for pkgName, pkg := range pkgmap {
		types[pkgName] = make(map[string]NamedType)

		for fname, f := range pkg.Files {
			fmt.Fprintln(os.Stderr, "File:", fname)

			ast.Inspect(f, func(n ast.Node) bool {
				switch nodeType := n.(type) {
				// skip comments
				case *ast.CommentGroup, *ast.Comment:
					return false
				case *ast.TypeSpec:
					types[pkgName][nodeType.Name.Name] = NamedType{
						Ident: nodeType.Name,
						Type:  nodeType.Type,
					}
					return false
				}

				return true
			})
		}

		// for n, _ := range pkg.Imports {
		// 	inspectDir(n)
		// }
	}

	return types
}
