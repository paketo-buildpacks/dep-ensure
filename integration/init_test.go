package integration_test

import (
	"os"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var settings struct {
	Buildpacks struct {
		DepEnsure struct {
			Online string
		}
	}
	Buildpack struct {
		ID   string
		Name string
	}
	// 	Config struct {
	// 		GoDist string `json:"go-dist"`
	// 	}
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	// file, err := os.Open("../integration.json")
	// Expect(err).NotTo(HaveOccurred())

	// Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())
	// Expect(file.Close()).To(Succeed())

	file, err := os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.DecodeReader(file, &settings)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	// root, err := filepath.Abs("./..")
	// Expect(err).ToNot(HaveOccurred())

	// buildpackStore := occam.NewBuildpackStore()

	// settings.Buildpacks.GoBuild.Online, err = buildpackStore.Get.
	// 	WithVersion("1.2.3").
	// 	Execute(root)
	// Expect(err).ToNot(HaveOccurred())

	// settings.Buildpacks.GoBuild.Offline, err = buildpackStore.Get.
	// 	WithVersion("1.2.3").
	// 	WithOfflineDependencies().
	// 	Execute(root)
	// Expect(err).ToNot(HaveOccurred())

	// settings.Buildpacks.GoDist.Online, err = buildpackStore.Get.
	// 	Execute(settings.Config.GoDist)
	// Expect(err).ToNot(HaveOccurred())

	// settings.Buildpacks.GoDist.Offline, err = buildpackStore.Get.
	// 	WithOfflineDependencies().
	// 	Execute(settings.Config.GoDist)
	// Expect(err).ToNot(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Default", testDefault)
	suite.Run(t)
}
