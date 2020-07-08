package depensure

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
)

//go:generate faux --interface BuildProcess --output fakes/build_process.go
type BuildProcess interface {
	Execute(workspace, goPath, gocachedir string) (err error)
}

func Build(
	buildProcess BuildProcess,
	logger LogEmitter,
	clock chronos.Clock,
) packit.BuildFunc {

	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		logger.Process("Executing build process")

		depcachedirLayer, err := context.Layers.Get("depcachedir", packit.CacheLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		gopath, err := ioutil.TempDir(os.TempDir(), "gopath")
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to create GOPATH dir: %w", err)
		}

		duration, err := clock.Measure(func() error {
			return buildProcess.Execute(context.WorkingDir, gopath, depcachedirLayer.Path)
		})

		if err != nil {
			return packit.BuildResult{}, err
		}

		err = os.RemoveAll(gopath)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to delete temp gopath dir: %w", err)
		}

		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		return packit.BuildResult{
			Layers:    []packit.Layer{depcachedirLayer},
			Processes: nil,
		}, nil
	}
}
