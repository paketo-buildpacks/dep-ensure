package depensure

import (
	// "github.com/paketo-buildpacks/packit/chronos"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/fs"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
)

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) (err error)
}

type DepEnsureProcess struct {
	executable Executable
	logs       scribe.Emitter
}

func NewDepEnsureProcess(executable Executable, logs scribe.Emitter) DepEnsureProcess {
	return DepEnsureProcess{
		executable: executable,
		logs:       logs,
	}
}

func (p DepEnsureProcess) Execute(workspace, gopath, depcachedir string) error {
	var err error
	tmpAppPath := filepath.Join(gopath, "src", "app")
	err = os.MkdirAll(tmpAppPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create GOPATH app dir: %w", err)
	}

	err = fs.Copy(workspace, tmpAppPath)
	if err != nil {
		return fmt.Errorf("failed to copy application source onto GOPATH: %w", err)
	}

	args := []string{"ensure"}
	buffer := bytes.NewBuffer(nil)

	p.logs.Subprocess("Running 'dep %s'", strings.Join(args, " "))
	err = p.executable.Execute(pexec.Execution{
		Args:   args,
		Dir:    tmpAppPath,
		Stdout: buffer,
		Stderr: buffer,
		Env:    append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath), fmt.Sprintf("DEPCACHEDIR=%s", depcachedir)),
	})

	if err != nil {
		p.logs.Detail(buffer.String())
		return fmt.Errorf("'dep ensure' command failed: %w", err)
	}

	err = os.RemoveAll(filepath.Join(workspace, "vendor"))
	if err != nil {
		return fmt.Errorf("failed to remove vendor from application source: %w", err)
	}

	err = fs.Copy(filepath.Join(tmpAppPath, "vendor"), filepath.Join(workspace, "vendor"))
	if err != nil {
		return fmt.Errorf("failed to copy vendor back to application source: %w", err)
	}

	err = fs.Copy(filepath.Join(tmpAppPath, "Gopkg.lock"), filepath.Join(workspace, "Gopkg.lock"))
	if err != nil {
		return fmt.Errorf("failed to copy Gopkg.lock back to application source: %w", err)
	}

	return err
}
