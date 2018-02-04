package libs

import "go/ast"

type NamedType struct {
	Ident *ast.Ident
	Type  ast.Expr
}
