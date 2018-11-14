package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/shims"
)

func main() {
	fmt.Printf("HEYYYYYY")
	buildpackDir := filepath.Join(os.Args[0], "..", "..")
	workspaceDir := filepath.Join(os.Args[1], "..")

	err := shims.Detect(&shims.Shim{}, buildpackDir, workspaceDir)
	if err != nil {
		log.Fatal(err)
	}
}
