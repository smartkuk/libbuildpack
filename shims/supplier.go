package shims

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Detector interface {
	Detect() error
}

type LifecycleBinaryRunner interface {
	RunBuild() error
}

type Supplier struct {
	Detector      Detector
	BinDir        string
	BuildDir      string
	BuildpacksDir string
	CacheDir      string
	DepsDir       string
	DepsIndex     string
	GroupMetadata string
	LaunchDir     string
	OrderMetadata string
	PlanMetadata  string
	WorkspaceDir  string
}

func (s *Supplier) Supply() error {
	if err := os.Symlink(s.BuildDir, filepath.Join(s.LaunchDir, "app")); err != nil {
		return err
	}

	if err := s.GetBuildPlan(); err != nil {
		return err
	}

	if err := s.RunLifeycleBuild(); err != nil {
		return err
	}

	if err := os.Remove(filepath.Join(s.LaunchDir, "app")); err != nil {
		return err
	}

	layers, err := filepath.Glob(filepath.Join(s.LaunchDir, "*"))
	if err != nil {
		return err
	}

	for _, layer := range layers {
		if filepath.Base(layer) == "config" {
			err = os.Rename(filepath.Join(s.LaunchDir, "config", "metadata.toml"), filepath.Join(s.BuildDir, "metadata.toml"))
			if err != nil {
				return err
			}
		} else {
			err := os.Rename(layer, filepath.Join(s.DepsDir, s.DepsIndex, filepath.Base(layer)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Supplier) GetBuildPlan() error {
	_, groupErr := os.Stat(s.GroupMetadata)
	_, planErr := os.Stat(s.PlanMetadata)

	if os.IsNotExist(groupErr) || os.IsNotExist(planErr) {
		if err := s.Detector.Detect(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Supplier) RunLifeycleBuild() error {
	cmd := exec.Command(
		filepath.Join(s.BinDir, "v3-builder"),
		"-buildpacks", s.BuildpacksDir,
		"-cache", s.CacheDir,
		"-group", s.GroupMetadata,
		"-launch", s.LaunchDir,
		"-plan", s.PlanMetadata,
		"-platform", s.WorkspaceDir,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PACK_STACK_ID=org.cloudfoundry.stacks."+os.Getenv("CF_STACK"))

	return cmd.Run()
}
