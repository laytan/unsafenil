package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/laytan/unsafenil/pkg/analyzer"
)

func TestNilNil(t *testing.T) {
	pkgs := []string{
		"examples",
	}
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), pkgs...)
}
