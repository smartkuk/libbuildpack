package shims

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

const setupPathContent = "export PATH={{ range $_, $path := . }}{{ $path }}:{{ end }}$PATH"

func Finalize(depsDir, depsIndex, profileDir string) error {
	files, err := filepath.Glob(filepath.Join(depsDir, depsIndex, "*", "*", "profile.d", "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Rename(file, filepath.Join(profileDir, filepath.Base(file)))
		if err != nil {
			return err
		}
	}

	BinDirs, err := filepath.Glob(filepath.Join(depsDir, depsIndex, "*", "*", "bin"))
	if err != nil {
		return err
	}

	for i, dir := range BinDirs {
		BinDirs[i] = strings.Replace(dir, filepath.Clean(depsDir), `$DEPS_DIR`, 1)
	}

	script, err := os.OpenFile(filepath.Join(profileDir, depsIndex+".sh"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer script.Close()

	setupPathTemplate, err := template.New("setupPathTemplate").Parse(setupPathContent)
	if err != nil {
		return err
	}

	return setupPathTemplate.Execute(script, BinDirs)
}

type inputMetadata struct {
	Processes []struct {
		Type    string
		Command string
	}
}

func (i *inputMetadata) findCommand(processType string) (string, error) {
	for _, p := range i.Processes {
		if p.Type == processType {
			return p.Command, nil
		}
	}
	return "", fmt.Errorf("unable to find process with type %s in launch metadata", processType)
}

type outputMetadata struct {
	DefaultProcessTypes struct {
		Web string
	} `yaml:"default_process_types"`
}

func Release(buildDir string, writer io.Writer) error {
	metadataFile, input := filepath.Join(buildDir, "metadata.toml"), inputMetadata{}
	_, err := toml.DecodeFile(metadataFile, &input)

	defer os.Remove(metadataFile)

	webCommand, err := input.findCommand("web")
	if err != nil {
		return err
	}

	output := outputMetadata{DefaultProcessTypes: struct{ Web string }{Web: webCommand}}
	return yaml.NewEncoder(writer).Encode(output)
}
