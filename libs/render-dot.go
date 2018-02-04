package libs

import (
	"os"
	"fmt"
	"bytes"
	"go/ast"
)

func RenderDot(out *os.File, pkgTypes map[string]map[string]NamedType) {
	fmt.Fprintf(out, "digraph %q { \n", "GoERD")

	var buf bytes.Buffer
	for pkg, types := range pkgTypes {

		fmt.Fprintf(out, "subgraph %q {\n", pkg)
		fmt.Fprintf(out, "label=%q;\n", pkg)

		// Nodes
		var i int
		for _, typ := range types {
			i++
			buf.Reset()

			switch t := typ.Type.(type) {
			case *ast.Ident:
				var label = fmt.Sprintf(`%s %s`, typ.Ident.Name, t.Name)
				fmt.Fprintf(out, " \"node-%s\" [shape=ellipse,label=\"%s\"];\n", typ.Ident.Name, escape(label))
			case *ast.SelectorExpr:
				var label = fmt.Sprintf(`%s %s`, typ.Ident.Name, toString(t))
				fmt.Fprintf(out, " \"node-%s\" [shape=ellipse,label=\"%s\"];\n", typ.Ident.Name, escape(label))
			case *ast.ChanType:
				var label = fmt.Sprintf(`%s %s`, typ.Ident.Name, toString(t))
				fmt.Fprintf(out, " \"node-%s\" [shape=box,label=\"%s\"];\n", typ.Ident.Name, escape(label))
			case *ast.FuncType:
				var label = fmt.Sprintf(`%s %s`, typ.Ident.Name, toString(t))
				fmt.Fprintf(out, " \"node-%s\" [shape=rectangle,label=\"%s\"];\n", typ.Ident.Name, escape(label))
			case *ast.ArrayType:
				var label = fmt.Sprintf(`%s %s`, typ.Ident.Name, toString(t))
				fmt.Fprintf(out, " \"node-%s\" [shape=rectangle,label=\"%s\"];\n", typ.Ident.Name, escape(label))
			case *ast.MapType:
				var label = fmt.Sprintf(`%s %s`, typ.Ident.Name, toString(t))
				fmt.Fprintf(out, " \"node-%snamedType\" [shape=rectangle,label=\"%s\"];\n", typ.Ident.Name, escape(label))
			case *ast.InterfaceType:
				fmt.Fprintf(&buf, `%s interface|`, typ.Ident.Name)
				for i, f := range t.Methods.List {
					if i > 0 {
						fmt.Fprintf(&buf, `|`)
					}
					fmt.Fprintf(&buf, `<f%d>`, i)
					// a,b,c Type
					for ii, n := range f.Names {
						fmt.Fprintf(&buf, "%s", n.Name)
						if ii > 0 {
							fmt.Fprintf(&buf, `,`)
						}
					}
					if len(f.Names) > 0 {
						fmt.Fprintf(&buf, ` `)
					}
					fmt.Fprintf(&buf, `%s`, toString(f.Type))
				}
				fmt.Fprintf(out, " \"node-%s\" [shape=Mrecord,label=\"{%s}\"];\n", typ.Ident.Name, escape(buf.String()))
			case *ast.StructType:
				fmt.Fprintf(&buf, `%s|`, typ.Ident.Name)
				for i, f := range t.Fields.List {
					if i > 0 {
						fmt.Fprintf(&buf, "|")
					}
					fmt.Fprintf(&buf, `<f%d>`, i)

					for ii, n := range f.Names {
						if ii > 0 {
							fmt.Fprintf(&buf, `, `)
						}
						fmt.Fprintf(&buf, `%s`, n.Name)
					}
					if len(f.Names) > 0 {
						fmt.Fprintf(&buf, ` `)
					}
					fmt.Fprintf(&buf, `%s`, toString(f.Type))
				}
				fmt.Fprintf(out, " \"node-%s\" [shape=record,label=\"{%s}\"];\n", typ.Ident.Name, escape(buf.String()))
			default:
				fmt.Fprintf(os.Stderr, "MISSED: %s: %#v\n ", toString(t), typ)
			}
		}

		// Edges
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

		fmt.Fprintf(out, "}\n")
	}
	fmt.Fprintf(out, "}\n\n")
}
