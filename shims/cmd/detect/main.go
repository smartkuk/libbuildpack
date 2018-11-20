package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/shims"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("incorrect number of arguments"))
	}

	buildpackDir, err := filepath.Abs(filepath.Join(os.Args[0], "..", ".."))
	if err != nil {
		log.Fatal(err)
	}
	workspaceDir, err := filepath.Abs(filepath.Join(os.Args[1], ".."))
	if err != nil {
		log.Fatal(err)
	}

	detector := shims.Detector{
		BinDir:        filepath.Join(buildpackDir, "bin"),
		BuildpacksDir: filepath.Join(buildpackDir, "cnbs"),
		GroupMetadata: filepath.Join(workspaceDir, "group.toml"),
		LaunchDir:     workspaceDir,
		OrderMetadata: filepath.Join(buildpackDir, "order.toml"),
		PlanMetadata:  filepath.Join(workspaceDir, "plan.toml"),
	}

	err = detector.Detect()
	if err != nil {
		log.Fatal(err)
	}
}
