package libs

import (
	"go/ast"
	"fmt"
	"os"
)

func RenderDotEdges(out *os.File, types map[string]NamedType) {
	for _, ptype := range types {
		switch t := ptype.Type.(type) {
		// TODO: exhaustive switch
		case *ast.FuncType:
			for i, typ := range dependsOn(t) {
				var from = fmt.Sprintf(`"node-%s":f%d`, ptype.Ident.Name, i)
				var to = fmt.Sprintf("node-%s", typ)
				if _, ok := types[typ]; ok {
					fmt.Fprintf(out, "%s -> %q;\n", from, to)
				}
			}
		case *ast.ChanType:
			for i, typ := range dependsOn(t) {
				var from = fmt.Sprintf(`"node-%s":f%d`, ptype.Ident.Name, i)
				var to = fmt.Sprintf("node-%s", typ)
				if _, ok := types[typ]; ok {
					fmt.Fprintf(out, "%s -> %q;\n", from, to)
				}
			}
		case *ast.InterfaceType:
			for i, f := range t.Methods.List {
				var from = fmt.Sprintf(`"node-%s":f%d`, ptype.Ident.Name, i)
				for _, typ := range dependsOn(f.Type) {
					var to = fmt.Sprintf("node-%s", typ)
					if _, ok := types[typ]; ok {
						fmt.Fprintf(out, "%s -> %q;\n", from, to)
					}
				}
			}
		case *ast.StructType:
			for i, f := range t.Fields.List {
				var from = fmt.Sprintf(`"node-%s":f%d`, ptype.Ident.Name, i)
				for _, typ := range dependsOn(f.Type) {
					var to = fmt.Sprintf("node-%s", typ)
					if _, ok := types[typ]; ok {
						fmt.Fprintf(out, "%s -> %q;\n", from, to)
					}
				}
			}
		}
	}
}
