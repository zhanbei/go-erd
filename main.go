package main

import (
	"flag"
	"os"
	"github.com/zhanbei/go-erd/libs"
)

//
// from: https://github.com/golang/example/tree/master/gotypes#typeandvalue
//
// running: go run ./cmd/goerd/main.go -path cmd/traverse/|dot -Tsvg > out.svg
func main() {
	var (
		path = flag.String("path", "", "path parse")
	)

	flag.Parse()

	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	libs.RenderDot(os.Stdout, libs.InspectDir(*path))
}
