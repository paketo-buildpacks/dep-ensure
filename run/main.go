package main

import (
	"os"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout)
	packit.Run(
		depensure.Detect(),
		depensure.Build(
			depensure.NewDepEnsureProcess(pexec.NewExecutable("dep"), logEmitter),
			logEmitter,
			chronos.DefaultClock,
		),
	)
}
