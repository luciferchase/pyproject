package script

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
)

//go:embed "samples/.gitignore"
var gitignore []byte

//go:embed "samples/.pre-commit-config.yaml"
var preCommitConfig []byte

//go:embed "samples/.pyproject.toml"
var pyprojectToml []byte

//go:embed "samples/.setup.cfg"
var setupCfg []byte

func addDeps(dir string, deps []string) error {
	var depFail string
	for _, dep := range deps {
		cmd := exec.Command("poetry", "add", "--dev", dep)
		if (dep == "black") {
			cmd.Args = append(cmd.Args, "--allow-prereleases")
		}
		cmd.Dir = path.Join(".", dir)

		if err := cmd.Run(); err != nil {
			depFail += dep + " "
		} else {
			log.Println("✔️ successfully added dependency:", dep)
		}
	}
	if depFail != "" {
		return errors.New("⚠ failed to add: " + depFail)
	}
	return nil
}

func addToFile(dir string) error {
	var fileFail string

	// Add .gitignore
	file, err := os.Create(path.Join(dir, ".gitignore"))
	if err != nil {
		log.Println("⚠️ failed to create .gitignore")
		fileFail += ".gitignore "
	}
	file.Write(gitignore)
	file.Close()

	// Add pre-commit config
	file, err = os.Create(path.Join(dir, "pre-commit-config.yaml"))
	if err != nil {
		log.Println("⚠️ failed to create pre-commit-config.yaml")
		fileFail += "pre-commit-config.yaml "
	}
	file.Write(preCommitConfig)
	file.Close()

	// Add to pyproject.toml
	file, err = os.OpenFile(
		path.Join(dir, "pyproject.toml"), os.O_APPEND, os.ModeAppend,
	)
	if err != nil {
		log.Println("⚠️ failed to open pyproject.toml")
		fileFail += "pyproject.toml "
	}
	file.Write(pyprojectToml)
	file.Close()

	// Add to setup.cfg
	file, err = os.Create(path.Join(dir, "setup.cfg"))
	if err != nil {
		log.Println("⚠️ failed to open setup.cfg")
		fileFail += "setup.cfg "
	}
	file.Write(setupCfg)
	file.Close()

	if fileFail != "" {
		return errors.New("⚠️ failed to add: " + fileFail)
	}
	return nil
}

func checkForGit() bool {
	if err := exec.Command("git", "--version").Run(); err != nil {
		return false
	}
	return true
}

func checkForPoetry() bool {
	if err := exec.Command("poetry", "--version").Run(); err != nil {
		return false
	}
	return true
}
