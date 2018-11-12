package main

import (
	"fmt"
	"github.com/cloudfoundry/libbuildpack/shims"
	"log"
	"os"
	"path/filepath"
)

func main() {
	buildpackDir := filepath.Join(os.Args[0], "..", "..")
	buildDir := os.Args[1]
	cacheDir := os.Args[2]
	depsDir := os.Args[3]
	depsIndex := os.Args[4]
	workspaceDir := filepath.Join(buildDir, "..")
	launchDir, err := filepath.Abs(filepath.Join("home", "vcap", "deps", depsIndex)) // figure this out (it's not home)
	fmt.Println("DIRS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("BUILDPACK DIR: " + buildpackDir)
	fmt.Println("build DIR: " + buildDir)
	fmt.Println("CACHE DIR: " + cacheDir)
	fmt.Println("deps DIR: " + depsDir)
	fmt.Println("deps index: " + depsIndex)
	fmt.Println("workspace DIR: " + workspaceDir)
	fmt.Println("launch DIR: " + launchDir)

	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(launchDir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(launchDir)


	err = shims.Supply(&shims.Shim{}, buildpackDir, buildDir, cacheDir, depsDir, depsIndex, workspaceDir, launchDir)
	if err != nil {
		log.Fatal(err)
	}
}