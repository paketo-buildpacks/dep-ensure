package depensure

import (
	"github.com/paketo-buildpacks/packit"
)

func Build(
	logs LogEmitter,
) packit.BuildFunc {

	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		return packit.BuildResult{
			Layers:    nil,
			Processes: nil,
		}, nil
	}
}
