// Command nonnil runs the nonnil analyzer standalone or as a `go vet` tool:
//
//	go install github.com/golang-cop/nonnil/cmd/nonnil@latest
//	go vet -vettool=$(which nonnil) ./...
package main

import (
	"github.com/golang-cop/nonnil"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nonnil.Analyzer) }
