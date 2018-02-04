package libs

import (
	"os"
	"fmt"
)

func RenderDot(out *os.File, pkgTypes map[string]map[string]NamedType) {
	fmt.Fprintf(out, "digraph %q { \n", "GoERD")

	for pkg, types := range pkgTypes {

		fmt.Fprintf(out, "subgraph %q {\n", pkg)
		fmt.Fprintf(out, "label=%q;\n", pkg)

		RenderDotNodes(out, types)
		RenderDotEdges(out, types)

		fmt.Fprintf(out, "}\n")
	}
	fmt.Fprintf(out, "}\n\n")
}
