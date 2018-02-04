package libs

import (
	"bytes"
	"go/ast"
	"fmt"
	"os"
)

func RenderDotNodes(out *os.File, types map[string]NamedType) {
	for _, typ := range types {
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
			var buf bytes.Buffer
			buf.Reset()
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
			var buf bytes.Buffer
			buf.Reset()
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
}
