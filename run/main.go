package main

import (
	"os"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	logEmitter := depensure.NewLogEmitter(os.Stdout)
	packit.Run(
		depensure.Detect(),
		depensure.Build(logEmitter),
	)
}
