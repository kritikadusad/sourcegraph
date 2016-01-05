package plan

import (
	"fmt"

	droneyaml "github.com/drone/drone-exec/yaml"
	"github.com/drone/drone/yaml/matrix"
	"src.sourcegraph.com/sourcegraph/pkg/inventory"
)

// autogenerateConfig consults a repo's inventory (of programming
// languages used) and generates a .drone.yml file that will build the
// repo. This is not guaranteed to be correct. The primary purpose of
// this auto-generated config is to fetch the project dependencies and
// compile the project to prepare for the srclib analysis step.
func autogenerateConfig(inv *inventory.Inventory) (*droneyaml.Config, []matrix.Axis, error) {
	// Merge the default configs for all of the languages we detect.
	var config droneyaml.Config
	matrix := matrix.Matrix{}
	for _, lang := range inv.Languages {
		c, ok := langConfigs[lang.Name]
		if !ok {
			c.build = buildLogMsg(fmt.Sprintf("automatic CI config does not yet support %s", lang.Name))
		}

		config.Build = append(config.Build, c.build)
		for key, vals := range c.matrix {
			matrix[key] = append(matrix[key], vals...)
		}
	}

	if len(config.Build) == 0 {
		config.Build = append(config.Build, buildLogMsg("no supported programming languages were auto-detected"))
	}

	return &config, calcMatrix(matrix), nil
}

var langConfigs = map[string]struct {
	build  droneyaml.BuildItem
	matrix map[string][]string
}{
	"Go": {
		build: droneyaml.BuildItem{
			Key: "Go $$GO_VERSION",
			Build: droneyaml.Build{
				Container: droneyaml.Container{Image: "golang:$$GO_VERSION"},
				Commands: []string{
					"go get -t ./...",
					"go build ./...",
					"go test -v ./...",
				},
			},
		},
		matrix: map[string][]string{"GO_VERSION": []string{"1.5"}},
	},
	"JavaScript": {
		build: droneyaml.BuildItem{
			Key: "JavaScript (node $$NODE_VERSION)",
			Build: droneyaml.Build{
				Container: droneyaml.Container{Image: "node:$$NODE_VERSION"},
				Commands: []string{
					"npm install",
					"npm run test",
				},
			},
		},
		matrix: map[string][]string{"NODE_VERSION": []string{"4"}},
	},
}
