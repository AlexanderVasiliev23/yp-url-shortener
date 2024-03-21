package osexitchecker

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func Test(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitAnalyzer, "./...")
}
