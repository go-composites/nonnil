// Command nonnil runs the nonnil analyzer standalone or as a `go vet` tool:
//
//	go install github.com/go-composites/nonnil/cmd/nonnil@latest
//	go vet -vettool=$(which nonnil) ./...
package main

import (
	"github.com/go-composites/nonnil"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nonnil.Analyzer) }
