package depensure

import (
	// "github.com/paketo-buildpacks/packit/chronos"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/fs"
	"github.com/paketo-buildpacks/packit/pexec"
)

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) (err error)
}

type DepEnsureProcess struct {
	executable Executable
	logs       LogEmitter
	// clock      chronos.Clock
}

func NewDepEnsureProcess(executable Executable, logs LogEmitter) DepEnsureProcess {
	return DepEnsureProcess{
		executable: executable,
		logs:       logs,
		// clock:      clock,
	}
}

func (p DepEnsureProcess) Execute(workspace, gopath string) error {
	var err error
	p.logs.Process("Executing build process")

	tmpAppPath := filepath.Join(gopath, "src", "app")
	err = os.MkdirAll(tmpAppPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create GOPATH app dir: %w", err)
	}

	err = fs.Copy(workspace, tmpAppPath)
	if err != nil {
		return fmt.Errorf("failed to copy application source onto GOPATH: %w", err)
	}

	buffer := bytes.NewBuffer(nil)
	err = p.executable.Execute(pexec.Execution{
		Args:   []string{"ensure"},
		Dir:    tmpAppPath,
		Stdout: buffer,
		Stderr: buffer,
		Env:    append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath)),
	})

	return err
}
