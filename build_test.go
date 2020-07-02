package depensure_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	depensure "github.com/paketo-buildpacks/dep-ensure"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir  string
		workingDir string
		cnbDir     string
		logs       *bytes.Buffer
		//timestamp  time.Time
		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = ioutil.TempDir("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		logs = bytes.NewBuffer(nil)
		build = depensure.Build(
			depensure.NewLogEmitter(logs),
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
			Layers:    nil,
			Processes: nil,
		}))

		Expect(logs.String()).To(ContainSubstring("Some Buildpack some-version"))
		Expect(logs.String()).To(ContainSubstring("Executing dep ensure"))
	})

	// context("when the workspace contents have not changed from a previous build", func() {
	// 	it.Before(func() {
	// 		layerContent := fmt.Sprintf("launch = true\n[metadata]\ncommand = \"some-start-command\"\nworkspace_sha = \"some-workspace-sha\"\nbuilt_at = %q\n", timestamp.Format(time.RFC3339Nano))
	// 		Expect(ioutil.WriteFile(filepath.Join(layersDir, "targets.toml"), []byte(layerContent), 0644)).To(Succeed())
	// 	})

	// 	it("skips the build process", func() {
	// 		result, err := build(packit.BuildContext{
	// 			WorkingDir: workingDir,
	// 			CNBPath:    cnbDir,
	// 			Stack:      "some-stack",
	// 			BuildpackInfo: packit.BuildpackInfo{
	// 				Name:    "Some Buildpack",
	// 				Version: "some-version",
	// 			},
	// 			Layers: packit.Layers{Path: layersDir},
	// 		})
	// 		Expect(err).NotTo(HaveOccurred())

	// 		Expect(result).To(Equal(packit.BuildResult{
	// 			Layers: []packit.Layer{
	// 				{
	// 					Name:      "targets",
	// 					Path:      filepath.Join(layersDir, "targets"),
	// 					SharedEnv: packit.Environment{},
	// 					BuildEnv:  packit.Environment{},
	// 					LaunchEnv: packit.Environment{},
	// 					Build:     false,
	// 					Launch:    true,
	// 					Cache:     false,
	// 					Metadata: map[string]interface{}{
	// 						"built_at":      timestamp.Format(time.RFC3339Nano),
	// 						"command":       "some-start-command",
	// 						"workspace_sha": "some-workspace-sha",
	// 					},
	// 				},
	// 				{
	// 					Name:      "gocache",
	// 					Path:      filepath.Join(layersDir, "gocache"),
	// 					SharedEnv: packit.Environment{},
	// 					BuildEnv:  packit.Environment{},
	// 					LaunchEnv: packit.Environment{},
	// 					Build:     false,
	// 					Launch:    false,
	// 					Cache:     true,
	// 				},
	// 			},
	// 			Processes: []packit.Process{
	// 				{
	// 					Type:    "web",
	// 					Command: "some-start-command",
	// 					Direct:  true,
	// 				},
	// 			},
	// 		}))

	// 		Expect(calculator.SumCall.Receives.Path).To(Equal(workingDir))
	// 		Expect(pathManager.SetupCall.CallCount).To(Equal(0))
	// 		Expect(buildProcess.ExecuteCall.CallCount).To(Equal(0))
	// 		Expect(pathManager.TeardownCall.CallCount).To(Equal(0))
	// 	})
	// })

	// context("failure cases", func() {
	// 	context("when the build process fails", func() {
	// 		it.Before(func() {
	// 			buildProcess.ExecuteCall.Returns.Err = errors.New("failed to execute build process")
	// 		})

	// 		it("returns an error", func() {
	// 			_, err := build(packit.BuildContext{
	// 				WorkingDir: workingDir,
	// 				CNBPath:    cnbDir,
	// 				Stack:      "some-stack",
	// 				BuildpackInfo: packit.BuildpackInfo{
	// 					Name:    "Some Buildpack",
	// 					Version: "some-version",
	// 				},
	// 				Layers: packit.Layers{Path: layersDir},
	// 			})
	// 			Expect(err).To(MatchError("failed to execute build process"))
	// 		})
	// 	})
	// })
}
