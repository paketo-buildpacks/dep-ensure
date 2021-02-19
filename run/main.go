package main

import (
	"os"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
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
