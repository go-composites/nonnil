package nonnil_test

import (
	"testing"

	"github.com/golang-cop/nonnil"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), nonnil.Analyzer, "a")
}
