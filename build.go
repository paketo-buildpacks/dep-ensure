package depensure

import (
	"fmt"
	"os"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface BuildProcess --output fakes/build_process.go
type BuildProcess interface {
	Execute(workspace, goPath, gocachedir string) (err error)
}

func Build(
	buildProcess BuildProcess,
	logger scribe.Emitter,
	clock chronos.Clock,
) packit.BuildFunc {

	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		logger.Break()

		logger.Process("Executing build process")

		depcachedirLayer, err := context.Layers.Get("depcachedir")
		if err != nil {
			return packit.BuildResult{}, err
		}

		depcachedirLayer.Cache = true

		gopath, err := os.MkdirTemp("", "gopath")
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to create GOPATH directory: %w", err)
		}

		duration, err := clock.Measure(func() error {
			return buildProcess.Execute(context.WorkingDir, gopath, depcachedirLayer.Path)
		})

		if err != nil {
			return packit.BuildResult{}, err
		}

		err = os.RemoveAll(gopath)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to delete GOPATH directory: %w", err)
		}

		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		return packit.BuildResult{
			Layers: []packit.Layer{depcachedirLayer},
		}, nil
	}
}
