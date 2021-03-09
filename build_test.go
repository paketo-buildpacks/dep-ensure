package depensure_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/dep-ensure/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir    string
		workingDir   string
		cnbDir       string
		buildProcess *fakes.BuildProcess
		logs         *bytes.Buffer
		timeStamp    time.Time
		clock        chronos.Clock
		build        packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = ioutil.TempDir("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		timeStamp = time.Now()
		clock = chronos.NewClock(func() time.Time {
			return timeStamp
		})

		buildProcess = &fakes.BuildProcess{}
		logs = bytes.NewBuffer(nil)
		build = depensure.Build(
			buildProcess,
			scribe.NewEmitter(logs),
			clock,
		)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a result that builds correctly", func() {
		result, err := build(packit.BuildContext{
			WorkingDir: workingDir,
			CNBPath:    cnbDir,
			Stack:      "some-stack",
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
			Layers: packit.Layers{Path: layersDir},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(result).To(Equal(packit.BuildResult{
			Layers: []packit.Layer{
				{
					Name:             "depcachedir",
					Path:             filepath.Join(layersDir, "depcachedir"),
					SharedEnv:        packit.Environment{},
					BuildEnv:         packit.Environment{},
					LaunchEnv:        packit.Environment{},
					ProcessLaunchEnv: map[string]packit.Environment{},
					Build:            false,
					Launch:           false,
					Cache:            true,
				},
			},
		}))

		Expect(buildProcess.ExecuteCall.CallCount).To(Equal(1))
		Expect(buildProcess.ExecuteCall.Receives.Workspace).To(Equal(workingDir))
		Expect(buildProcess.ExecuteCall.Returns.Err).To(BeNil())
		Expect(logs.String()).To(ContainSubstring("Some Buildpack some-version"))
		Expect(logs.String()).To(ContainSubstring("Executing build process"))
		Expect(logs.String()).To(ContainSubstring("Completed in "))
	})

	context("failure cases", func() {
		context("when the get for the depcachedir fails", func() {
			it.Before(func() {
				Expect(os.Chmod(layersDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(layersDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					Stack:      "some-stack",
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Layers: packit.Layers{Path: layersDir},
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when the build process fails", func() {
			it.Before(func() {
				buildProcess.ExecuteCall.Returns.Err = errors.New("failed to execute build process")
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					CNBPath:    cnbDir,
					Stack:      "some-stack",
					BuildpackInfo: packit.BuildpackInfo{
						Name:    "Some Buildpack",
						Version: "some-version",
					},
					Layers: packit.Layers{Path: layersDir},
				})
				Expect(err).To(MatchError("failed to execute build process"))
			})
		})
	})
}
