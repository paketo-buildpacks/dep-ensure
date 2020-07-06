package depensure_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitGoBuild(t *testing.T) {
	suite := spec.New("dep-ensure", spec.Report(report.Terminal{}))
	// suite("Build", testBuild)
	// suite("Detect", testDetect)
	suite("DepEnsureProcess", testDepEnsureProcess)
	suite.Run(t)
}
