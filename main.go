package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/laytan/unsafenil/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
