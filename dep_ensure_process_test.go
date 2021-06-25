package depensure_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/dep-ensure/fakes"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDepEnsureProcess(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workspace     string
		gopath        string
		depcachedir   string
		executable    *fakes.Executable
		buffer        *bytes.Buffer
		commandOutput *bytes.Buffer

		process depensure.DepEnsureProcess
	)

	it.Before(func() {
		var err error

		workspace, err = os.MkdirTemp("", "workspace")
		Expect(err).NotTo(HaveOccurred())

		gopath, err = os.MkdirTemp("", "gopath")
		Expect(err).NotTo(HaveOccurred())

		depcachedir, err = os.MkdirTemp("", "depcachedir")
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(workspace, "test.go"), nil, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = os.MkdirAll(filepath.Join(workspace, "dir1", "dir2"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(workspace, "dir1", "dir2", "somefile.go"), nil, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		executable = &fakes.Executable{}
		executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
			err = os.MkdirAll(filepath.Join(gopath, "src", "app", "vendor", "somedir"), os.ModePerm)
			if err != nil {
				return err
			}

			if err := os.WriteFile(filepath.Join(gopath, "src", "app", "Gopkg.lock"), nil, 0755); err != nil {
				return err
			}

			return nil
		}

		buffer = bytes.NewBuffer(nil)
		commandOutput = bytes.NewBuffer(nil)

		process = depensure.NewDepEnsureProcess(executable, scribe.NewEmitter(buffer))
	})

	it.After(func() {
		Expect(os.RemoveAll(workspace)).To(Succeed())
		Expect(os.RemoveAll(gopath)).To(Succeed())
	})

	context("Execute", func() {
		it("succeeds", func() {
			Expect(process.Execute(workspace, gopath, depcachedir)).To(Succeed())
			Expect(executable.ExecuteCall.Receives.Execution).To(Equal(pexec.Execution{
				Args:   []string{"ensure"},
				Dir:    filepath.Join(gopath, "src", "app"),
				Stdout: commandOutput,
				Stderr: commandOutput,
				Env:    append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath), fmt.Sprintf("DEPCACHEDIR=%s", depcachedir)),
			}))

			Expect(filepath.Join(gopath, "src", "app", "test.go")).To(BeAnExistingFile())

			// make sure the file moves do not mess with src files
			Expect(filepath.Join(workspace, "dir1", "dir2", "somefile.go")).To(BeAnExistingFile())
			Expect(filepath.Join(workspace, "Gopkg.lock")).To(BeAnExistingFile())
			Expect(filepath.Join(workspace, "vendor", "somedir")).To(BeADirectory())

			Expect(buffer.String()).To(ContainSubstring("    Running 'dep ensure'"))
		})

		context("failure cases", func() {
			context("when unable to write to the tmp gopath dir", func() {
				it.Before(func() {
					Expect(os.Chmod(gopath, 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(gopath, 0777)).To(Succeed())
				})

				it("returns an error", func() {
					err := process.Execute(workspace, gopath, depcachedir)
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})

			context("when unable to write vendor to workspace dir", func() {
				it.Before(func() {
					Expect(os.Chmod(workspace, 0555)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(workspace, 0777)).To(Succeed())
				})

				it("returns an error", func() {
					err := process.Execute(workspace, gopath, depcachedir)
					Expect(err).To(MatchError(ContainSubstring("failed to copy vendor")))
				})
			})

			context("when the executable fails", func() {
				it.Before(func() {
					executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
						fmt.Fprintln(execution.Stdout, "dep ensure error on stdout")
						fmt.Fprintln(execution.Stderr, "dep ensure error on stderr")
						return errors.New("failed to execute")
					}
				})

				it("returns an error", func() {
					err := process.Execute(workspace, gopath, depcachedir)
					Expect(buffer.String()).To(ContainSubstring("dep ensure error on stdout\n"))
					Expect(buffer.String()).To(ContainSubstring("dep ensure error on stderr\n"))
					Expect(err).To(MatchError("'dep ensure' command failed: failed to execute"))
				})
			})
		})
	})
}
