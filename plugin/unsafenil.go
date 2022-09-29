package main

import (
	"github.com/laytan/unsafenil/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

var AnalyzerPlugin = analyzerPlugin{}

type analyzerPlugin struct{}

func (a *analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.New()}
}
