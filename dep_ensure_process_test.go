package depensure_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/dep-ensure/fakes"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDepEnsureProcess(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workspace     string
		gopath        string
		executable    *fakes.Executable
		buffer        *bytes.Buffer
		commandOutput *bytes.Buffer

		process depensure.DepEnsureProcess
	)

	it.Before(func() {
		var err error

		workspace, err = ioutil.TempDir("", "workspace")
		Expect(err).NotTo(HaveOccurred())

		gopath, err = ioutil.TempDir("", "gopath")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workspace, "test.go"), nil, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = os.MkdirAll(filepath.Join(workspace, "dir1", "dir2"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workspace, "dir1", "dir2", "somefile.go"), nil, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		executable = &fakes.Executable{}
		executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
			err = os.MkdirAll(filepath.Join(gopath, "src", "app", "vendor"), os.ModePerm)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(filepath.Join(gopath, "src", "app", "Gopkg.lock"), nil, 0755); err != nil {
				return err
			}

			return nil
		}

		buffer = bytes.NewBuffer(nil)
		commandOutput = bytes.NewBuffer(nil)

		process = depensure.NewDepEnsureProcess(executable, depensure.NewLogEmitter(buffer))
	})

	it.After(func() {
		Expect(os.RemoveAll(workspace)).To(Succeed())
		Expect(os.RemoveAll(gopath)).To(Succeed())
	})

	context("Execute", func() {
		it("succeeds", func() {
			var err error

			Expect(process.Execute(workspace, gopath)).To(Succeed())
			Expect(executable.ExecuteCall.Receives.Execution).To(Equal(pexec.Execution{
				Args:   []string{"ensure"},
				Dir:    filepath.Join(gopath, "src", "app"),
				Stdout: commandOutput,
				Stderr: commandOutput,
				Env:    append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath)),
			}))

			_, err = os.Stat(filepath.Join(gopath, "src", "app", "test.go"))
			Expect(err).NotTo(HaveOccurred())

			// make sure the file moves do not mess with src files
			_, err = os.Stat(filepath.Join(workspace, "dir1", "dir2", "somefile.go"))
			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(filepath.Join(workspace, "Gopkg.lock"))
			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(filepath.Join(workspace, "vendor"))
			Expect(err).NotTo(HaveOccurred())

			Expect(buffer.String()).To(ContainSubstring("    Running 'dep ensure'"))
		})

		context("failure cases", func() {
			context("when unable to write to workspace dir", func() {
				it.Before(func() {
					Expect(os.Chmod(workspace, 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(workspace, 0777)).To(Succeed())
				})

				it("returns an error", func() {
					err := process.Execute(workspace, gopath)
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
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
					err := process.Execute(workspace, gopath)
					Expect(buffer.String()).To(ContainSubstring("dep ensure error on stdout\n"))
					Expect(buffer.String()).To(ContainSubstring("dep ensure error on stderr\n"))
					Expect(err).To(MatchError("'dep ensure' command failed: failed to execute"))
				})
			})
		})
	})
}
