/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package watchexec

import (
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		build Build
		ctx   libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = os.MkdirTemp("", "build")
		Expect(err).NotTo(HaveOccurred())

		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "watchexec"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "watchexec",
					"version": "1.15.3",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
	})

	it("contributes Watchexec for API <= 0.6", func() {
		ctx.Buildpack.API = "0.6"

		result, err := build.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("watchexec"))

		Expect(result.BOM.Entries).To(HaveLen(1))
		Expect(result.BOM.Entries[0].Name).To(Equal("watchexec"))
	})

	it("contributes Watchexec for API 0.7+", func() {
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "watchexec",
					"version": "1.15.3",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []string{"cpe:2.3:a:watchexec:watchexec:1.17.1:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/watchexec@1.17.1?arch=amd64",
				},
			},
		}
		result, err := build.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("watchexec"))

		Expect(result.BOM.Entries).To(HaveLen(1))
		Expect(result.BOM.Entries[0].Name).To(Equal("watchexec"))
	})
}
