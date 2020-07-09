package main

import (
	"os"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
)

func main() {
	logEmitter := depensure.NewLogEmitter(os.Stdout)
	packit.Run(
		depensure.Detect(),
		depensure.Build(
			depensure.NewDepEnsureProcess(pexec.NewExecutable("dep"), logEmitter),
			logEmitter,
			chronos.DefaultClock,
		),
	)
}
