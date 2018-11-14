package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/buildpack/lifecycle"
	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	// buildpackDir := filepath.Join(os.Args[0], "..", "..")
	// buildDir := os.Args[1]
	// cacheDir := os.Args[2]
	// depsDir := os.Args[3]
	// depsIndex := os.Args[4]
	// workspaceDir := filepath.Join(buildDir, "..")
	// launchDir := filepath.Join(string(filepath.Separator), "home", "vcap", "deps", depsIndex)

	// err := os.MkdirAll(launchDir, 0777)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer os.RemoveAll(launchDir)

	// err = shims.Supply(&shims.Shim{}, buildpackDir, buildDir, cacheDir, depsDir, depsIndex, workspaceDir, launchDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// buildpackDir := os.Args[0]

	// buildpacks, err := lifecycle.NewBuildpackMap(buildpackDir)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	// How do we get access to buildpack map instead of hardcoding it here?
	bpMap := lifecycle.BuildpackMap{
		"org.cloudfoundry.buildpacks.npm@latest": {
			Name: "npm",
		},
		"org.cloudfoundry.buildpacks.nodejs@latest": {
			Name: "nodejs",
		},
	}

	var order lifecycle.BuildpackOrder
	orderPath, _ := filepath.Abs(filepath.Join("order.toml"))
	order, err := bpMap.ReadOrder(orderPath)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("\n\nbuildpack name %v\n", order[0].Buildpacks[0].Name)
	fmt.Printf("buildpack version %v\n", order[0].Buildpacks[0].Version)
	// group, err := lifecycle.ReadOrderTOML(orderPath)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	// Search manifest for bp name +version for URl
	manifestDir := filepath.Join("manifest.yml")
	manifest, err = libbuildpack.NewManifest(manifestDir, nil, time.Now())
	// Install node-v3
	// Install npm v3

}
