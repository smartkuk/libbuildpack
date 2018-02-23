package packager

//go:generate go-bindata -pkg $GOPACKAGE -prefix scaffold scaffold/...

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudfoundry/libbuildpack"
	yaml "gopkg.in/yaml.v2"
)

var CacheDir = filepath.Join(os.Getenv("HOME"), ".buildpack-packager", "cache")
var Stdout, Stderr io.Writer = os.Stdout, os.Stderr

func CompileExtensionPackage(bpDir, version string, cached bool) (string, error) {
	bpDir, err := filepath.Abs(bpDir)
	if err != nil {
		return "", err
	}
	dir, err := copyDirectory(bpDir)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(dir, "VERSION"), []byte(version), 0644)
	if err != nil {
		return "", err
	}

	isCached := "--uncached"
	if cached {
		isCached = "--cached"
	}
	cmd := exec.Command("bundle", "exec", "buildpack-packager", isCached)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	cmd.Env = append(os.Environ(), "BUNDLE_GEMFILE=cf.Gemfile")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var manifest struct {
		Language string `yaml:"language"`
	}

	if err := libbuildpack.NewYAML().Load(filepath.Join(bpDir, "manifest.yml"), &manifest); err != nil {
		return "", err
	}

	zipFile := fmt.Sprintf("%s_buildpack-v%s.zip", manifest.Language, version)
	if cached {
		zipFile = fmt.Sprintf("%s_buildpack-cached-v%s.zip", manifest.Language, version)
	}
	if err := libbuildpack.CopyFile(filepath.Join(dir, zipFile), filepath.Join(bpDir, zipFile)); err != nil {
		return "", err
	}

	return filepath.Join(dir, zipFile), nil
}

func checksumHex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// START AUTO-GENERATED CODE (from bindata.go)
// RestoreAsset restores an asset under the given directory
// TODO: maybe make scaffold into a struct and split up packager and scaffold if they don't share anything
func OurRestoreAsset(dir, name string, funcMap template.FuncMap, shas map[string]string) error {
	data, err := Asset(name)

	if err != nil {
		return err
	}

	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	// START NON-AUTO-GENERATED CODE
	t, err := template.New("").Funcs(funcMap).Parse(string(data))
	if err != nil {
		return err
	}
	var b bytes.Buffer
	f := bufio.NewWriter(&b)
	if err := t.Execute(f, nil); err != nil {
		return err
	}
	if err := f.Flush(); err != nil {
		return err
	}

	if strings.HasPrefix(name, "src/LANGUAGE/") {
		langName, ok := funcMap["LANGUAGE"].(func() string)
		if !ok {
			return fmt.Errorf("Could not find language from funcmap")
		}
		name = strings.Replace(name, "src/LANGUAGE/", fmt.Sprintf("src/%s/", langName()), 1)
	}

	oldSha256 := ""
	if oldContents, err := ioutil.ReadFile(_filePath(dir, name)); err == nil {
		oldSha256 = checksumHex(oldContents)
	}

	if name == "bin/supply" {
		fmt.Println(name, oldSha256, shas[name])
	}
	if shas[name] != "" && shas[name] != oldSha256 {
		fmt.Fprintf(Stdout, "***Ignoring %s becuase it has been modified***\n", name)
		return nil
	}
	shas[name] = checksumHex(b.Bytes())

	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(_filePath(dir, name), b.Bytes(), info.Mode()); err != nil {
		return err
	}
	// END NON-AUTO-GENERATED CODE

	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}

	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func OurRestoreAssets(dir, name string, funcMap template.FuncMap, shas map[string]string) error {
	//TODO: is passing in the shas map the best way to go about this?
	children, err := AssetDir(name)
	// File
	if err != nil {
		err = OurRestoreAsset(dir, name, funcMap, shas)
		if err != nil {
			return err
		}
	}

	// Dir
	for _, child := range children {
		err = OurRestoreAssets(dir, filepath.Join(name, child), funcMap, shas)
		if err != nil {
			return err
		}
	}
	return nil
}

// END AUTO-GENERATED CODE (from bindata.go)

func Scaffold(bpDir string, languageName string) error {
	language := func() string {
		return languageName
	}

	funcMap := template.FuncMap{
		"LANGUAGE": language,
	}

	type sha struct {
		Sha map[string]string `yaml:"sha"`
	}
	var shas = map[string]string{}
	if found, err := libbuildpack.FileExists(filepath.Join(bpDir, "sha.yml")); err != nil {
		return err
	} else if found {
		shas, err = readShaYML(bpDir)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(Stdout, "Creating directory and files")
	if err := OurRestoreAssets(bpDir, "", funcMap, shas); err != nil {
		return err
	}

	libbuildpack.NewYAML().Write(filepath.Join(bpDir, "sha.yml"), sha{
		Sha: shas,
	})

	if false {
		// Install dep and download dependencies (gomega, ginkgo, libbuildpack, etc)
		fmt.Fprintln(Stdout, "Installing dep")
		cmd := exec.Command("go", "get", "-u", "github.com/golang/dep/cmd/dep")
		cmd.Stdout = Stdout
		cmd.Stderr = Stderr
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOBIN=%s/.bin", bpDir), fmt.Sprintf("GOPATH=%s", bpDir))
		cmd.Dir = bpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go get -u github.com/golang/dep/cmd/dep: %s", err)
		}
		// TODO delete or uncommment
		// if err := os.Rename(filepath.Join(bpDir, "src", "github.com"), filepath.Join(bpDir, "src", languageName, "vendor", "github.com")); err != nil {
		// 	return err
		// }
		fmt.Fprintln(Stdout, "Running dep ensure")
		cmd = exec.Command(filepath.Join(bpDir, ".bin", "dep"), "ensure")
		cmd.Stdout = Stdout
		cmd.Stderr = Stderr
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOBIN=%s/.bin", bpDir), fmt.Sprintf("GOPATH=%s", bpDir))
		cmd.Dir = filepath.Join(bpDir, "src", languageName)
		if err := cmd.Run(); err != nil {
			fmt.Printf("GOPATH=%s\n", bpDir)
			return fmt.Errorf("dep ensure: %s", err)
		}
	}

	return nil
}

func Upgrade(bpDir string) error {
	manifest, err := readManifest(bpDir)
	if err != nil {
		return fmt.Errorf("error opening manifest: %s", err)
	}

	return Scaffold(bpDir, manifest.Language)
}

func readShaYML(bpDir string) (map[string]string, error) {
	type sha struct {
		Sha map[string]string `yaml:"sha"`
	}
	shas := &sha{}
	data, err := ioutil.ReadFile(filepath.Join(bpDir, "sha.yml"))
	if err != nil {
		return map[string]string{}, err
	}
	if err := yaml.Unmarshal(data, shas); err != nil {
		return map[string]string{}, err
	}

	return shas.Sha, nil
}

func readManifest(bpDir string) (*Manifest, error) {
	manifest := &Manifest{}
	data, err := ioutil.ReadFile(filepath.Join(bpDir, "manifest.yml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func Package(bpDir, cacheDir, version string, cached bool) (string, error) {
	bpDir, err := filepath.Abs(bpDir)
	if err != nil {
		return "", err
	}
	dir, err := copyDirectory(bpDir)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(dir, "VERSION"), []byte(version), 0644)
	if err != nil {
		return "", err
	}

	manifest, err := readManifest(dir)
	if err != nil {
		return "", err
	}

	if manifest.PrePackage != "" {
		cmd := exec.Command(manifest.PrePackage)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(Stdout, string(out))
			return "", err
		}
	}

	files := []File{}
	for _, name := range manifest.IncludeFiles {
		files = append(files, File{name, filepath.Join(dir, name)})
	}

	if cached {
		var m map[string]interface{}
		if err := libbuildpack.NewYAML().Load(filepath.Join(dir, "manifest.yml"), &m); err != nil {
			return "", err
		}
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			log.Fatalf("error: %v", err)
		}
		for idx, d := range manifest.Dependencies {
			file := filepath.Join("dependencies", fmt.Sprintf("%x", md5.Sum([]byte(d.URI))), filepath.Base(d.URI))
			if err := setFileOnDep(m, idx, file); err != nil {
				return "", err
			}

			if _, err := os.Stat(filepath.Join(cacheDir, file)); err != nil {
				if err := downloadFromURI(d.URI, filepath.Join(cacheDir, file)); err != nil {
					return "", err
				}
			}

			if err := checkSha256(filepath.Join(cacheDir, file), d.SHA256); err != nil {
				return "", err
			}

			files = append(files, File{file, filepath.Join(cacheDir, file)})
		}
		if err := libbuildpack.NewYAML().Write(filepath.Join(dir, "manifest.yml"), m); err != nil {
			return "", err
		}
	}

	zipFile := fmt.Sprintf("%s_buildpack-v%s.zip", manifest.Language, version)
	if cached {
		zipFile = fmt.Sprintf("%s_buildpack-cached-v%s.zip", manifest.Language, version)
	}
	zipFile = filepath.Join(bpDir, zipFile)

	ZipFiles(zipFile, files)

	return zipFile, err
}

func setFileOnDep(m map[string]interface{}, idx int, file string) error {
	if deps, ok := m["dependencies"].([]interface{}); ok {
		if dep, ok := deps[idx].(map[interface{}]interface{}); ok {
			dep["file"] = file
		} else {
			return fmt.Errorf("Could not cast deps[idx] to map[interface{}]interface{}")
		}
	} else {
		return fmt.Errorf("Could not cast dependencies to []interface{}")
	}
	return nil
}

func downloadFromURI(uri, fileName string) error {
	err := os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	output, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer output.Close()

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	var source io.ReadCloser

	if u.Scheme == "file" {
		source, err = os.Open(u.Path)
		if err != nil {
			return err
		}
		defer source.Close()
	} else {
		response, err := http.Get(uri)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		source = response.Body

		if response.StatusCode < 200 || response.StatusCode > 299 {
			return fmt.Errorf("could not download: %d", response.StatusCode)
		}
	}

	_, err = io.Copy(output, source)

	return err
}

func checkSha256(filePath, expectedSha256 string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	sum := sha256.Sum256(content)

	actualSha256 := hex.EncodeToString(sum[:])

	if actualSha256 != expectedSha256 {
		return fmt.Errorf("dependency sha256 mismatch: expected sha256 %s, actual sha256 %s", expectedSha256, actualSha256)
	}
	return nil
}

func ZipFiles(filename string, files []File) error {
	newfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file.Path)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate
		header.Name = file.Name

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyDirectory(srcDir string) (string, error) {
	destDir, err := ioutil.TempDir("", "buildpack-packager")
	if err != nil {
		return "", err
	}
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		path, err = filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		if path == ".git" || path == "tests" {
			return filepath.SkipDir
		}

		dest := filepath.Join(destDir, path)
		if info.IsDir() {
			err = os.MkdirAll(dest, info.Mode())
			if err != nil {
				return err
			}
		} else {
			src, err := os.Open(filepath.Join(srcDir, path))
			if err != nil {
				return err
			}
			defer src.Close()

			err = os.MkdirAll(filepath.Dir(dest), 0755)
			if err != nil {
				return err
			}

			fh, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(fh, src)
			fh.Close()
			return err
		}
		return nil
	})
	return destDir, err
}
