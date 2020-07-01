package main

import (
	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	packit.Run(
		depensure.Detect(),
		depensure.Build(),
	)
}
