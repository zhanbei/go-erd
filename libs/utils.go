package libs

import (
	"strings"
	"go/ast"
	"bytes"
	"fmt"
)

func escape(s string) string {
	for _, ch := range " '`[]{}()*" {
		s = strings.Replace(s, string(ch), `\`+string(ch), -1)
	}

	return s
}

func toString(n interface{}) string {
	switch t := n.(type) {
	case nil:
		return "nil"
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return toString(t.X) + "." + toString(t.Sel)
	case *ast.Object:
		return t.Name
	case *ast.StarExpr:
		return `*` + toString(t.X)
	case *ast.InterfaceType:
		// TODO:
		return `interface{}`
	case *ast.MapType:
		return `map[` + toString(t.Key) + `]` + toString(t.Value)
	case *ast.ChanType:
		return `chan ` + toString(t.Value)
	case *ast.StructType:
		// TODO:
		return `struct {}` //+ toString(t.)
	case *ast.Ellipsis:
		return `...` + toString(t.Elt)
	case *ast.Field:
		// ignoring names
		return toString(t.Type)

	case *ast.FuncType:
		var buf bytes.Buffer
		fmt.Fprint(&buf, `func(`)
		if t.Params != nil && len(t.Params.List) > 0 {
			for i, p := range t.Params.List {
				if i > 0 {
					fmt.Fprint(&buf, `, `)
				}
				fmt.Fprint(&buf, toString(p))
			}
		}
		fmt.Fprint(&buf, `)`)

		if t.Results != nil && len(t.Results.List) > 0 {
			fmt.Fprint(&buf, ` (`)
			for i, r := range t.Results.List {
				if i > 0 {
					fmt.Fprint(&buf, `, `)
				}
				fmt.Fprint(&buf, toString(r))
			}
			fmt.Fprint(&buf, `)`)
		}

		return buf.String()
	case *ast.ArrayType:
		return `[]` + toString(t.Elt)
	default:
		return fmt.Sprintf("%#v", n)
	}
}

// collect all the type names node n depends on
func dependsOn(n interface{}) []string {
	switch t := n.(type) {
	case nil:
		return nil
	case *ast.Ident:
		return []string{t.Name}
	case *ast.SelectorExpr:
		return []string{toString(t.X) + "." + t.Sel.Name}
	case *ast.Object:
		return []string{t.Name}
	case *ast.Field:
		return dependsOn(t.Type)
	case *ast.StarExpr:
		return dependsOn(t.X)
	case *ast.MapType:
		return append(dependsOn(t.Key), dependsOn(t.Value)...)
	case *ast.ChanType:
		return dependsOn(t.Value)
	case *ast.InterfaceType:
		if t.Methods == nil {
			return nil
		}
		var types []string
		for _, v := range t.Methods.List {
			types = append(types, dependsOn(v.Type)...)
		}
		return types
	case *ast.StructType:
		var types []string
		for _, v := range t.Fields.List {
			types = append(types, dependsOn(v.Type)...)
		}
		return types
	case *ast.FuncType:
		var types []string

		if t.Params != nil {
			for _, v := range t.Params.List {
				types = append(types, dependsOn(v.Type)...)
			}
		}

		if t.Results != nil {
			for _, v := range t.Results.List {
				types = append(types, dependsOn(v.Type)...)
			}
		}

		return types

	case *ast.ArrayType:
		return dependsOn(t.Elt)
	default:
		return []string{fmt.Sprintf("%#v", n)}
	}
}
