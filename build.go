package depensure

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface BuildProcess --output fakes/build_process.go
type BuildProcess interface {
	Execute(workspace, goPath string) (err error)
}

func Build(
	buildProcess BuildProcess,
	logger LogEmitter,
) packit.BuildFunc {

	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		logger.Process("Executing build process")

		gopath, err := ioutil.TempDir(os.TempDir(), "gopath")
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to create GOPATH dir: %w", err)
		}

		err = buildProcess.Execute(context.WorkingDir, gopath)
		if err != nil {
			return packit.BuildResult{}, err
		}

		return packit.BuildResult{
			Layers:    nil,
			Processes: nil,
		}, nil
	}
}
