package watchexec

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/libpak"
)

func testWatchexec(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path, err = ioutil.TempDir("", "watchexec-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Watchexec", func() {
		dep := libpak.BuildpackDependency{
			URI:    "http://localhost:8080/stub-watchexec.tar.xz",
			SHA256: "e7ef3171f63f7568bb6e2dbb5e2ca0cf83ece8330b0a72a6a751b0dc4f5d4ac5",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		j, _ := NewWatchexec(dep, dc)
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.LayerTypes.Build).To(BeFalse())
		Expect(layer.LayerTypes.Cache).To(BeFalse())
		Expect(layer.LayerTypes.Launch).To(BeTrue())

		Expect(filepath.Join(layer.Path, "watchexec")).To(BeARegularFile())
	})

}
