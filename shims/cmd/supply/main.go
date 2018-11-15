package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry/libbuildpack"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	type buildpack struct {
		Id string
		Version string
	}
	type group struct {
		Labels []string
		Buildpacks []buildpack
	}
	type Order struct {
		Groups []group
	}

	cwd, _ := os.Getwd()
	orderPath := filepath.Join(cwd, "order.toml")
	fmt.Println("orderpath", orderPath)

	var order Order
	_, err := toml.DecodeFile(orderPath, &order)
	if err != nil {
		fmt.Println(err)
	}

	manifest, err := libbuildpack.NewManifest(cwd, libbuildpack.NewLogger(os.Stdout), time.Now())
	for _, group := range order.Groups {
		for _, bp := range group.Buildpacks {

			ids := strings.Split(bp.Id, ".")
			name := fmt.Sprintf("%s-cnb", ids[len(ids)-1])
			dep, err := manifest.DefaultVersion(name)
			if err != nil {
				fmt.Println(err)
			}
			entry, err := manifest.GetEntry(dep)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("entry %+v\n", entry)
		}
	}

}
