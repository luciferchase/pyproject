package script

import (
	"errors"
	"log"
	"os/exec"
	"path"
)

var dependencies = []string{
	"black", "flake8", "isort", "mypy", "pre-commit", "pytest-cov",
}

func Run(name string) error {
	// init a project
	log.Println("initializing project...")
	cmd := exec.Command("poetry", "new", name)
	if err := cmd.Run(); err != nil {
		if checkForPoetry() {
			return errors.New("🛑 poetry new failed")
		} else {
			log.Println("🛑 poetry not installed: please install poetry",
			"https://python-poetry.org/docs/")
			return errors.New("🛑 poetry not installed")
		}
	}

	log.Println("✅ successfully initialized project")
	exec.Command("cd", name).Run()

	// Add bunch of dev dependencies
	log.Println("adding dependencies may take a while...")
	if err := addDeps(name, dependencies); err != nil {
		return err
	}
	log.Println("✅ all dependencies added successfully")

	// Add config files
	if err := addToFile(name); err != nil {
		return err
	}
	log.Println("✅ config files added successfully")

	// Initialize git
	log.Println("initializing git...")
	cmd = exec.Command("git", "init", "-b", "main")
	cmd.Dir = path.Join(".", name)
	if err := cmd.Run(); err != nil {
		if checkForGit() {
			return errors.New("🛑 git init failed")
		} else {
			log.Println("🛑 git not installed: please install git",
			"https://git-scm.com/")
			return errors.New("🛑 git not installed")
		}
	}
	log.Println("✅ git initialized successfully")

	// initialize git pre-commit hook
	log.Println("initializing git pre-commit hook...")

	cmd = exec.Command("poetry", "run", "pre-commit", "install")
	cmd.Dir = path.Join(".", name)
	if err := cmd.Run(); err != nil {
		return errors.New("⚠️ poetry pre-commit install failed")
	}

	cmd = exec.Command("poetry", "run", "pre-commit", "autoupdate")
	cmd.Dir = path.Join(".", name)
	if err := cmd.Run(); err != nil {
		return errors.New("⚠️ poetry pre-commit autoupdate failed")
	}
	log.Println("✅ poetry pre-commit hook installed successfully")

	// Add first commit
	log.Println("adding everything to git...")

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = path.Join(".", name)
	if err := cmd.Run(); err != nil {
		return errors.New("⚠️ git add failed")
	}

	cmd = exec.Command("git", "commit", "-m", "🚀 initial commit")
	cmd.Dir = path.Join(".", name)
	if err := cmd.Run(); err != nil {
		return errors.New("⚠ git commit failed: " +
			"can be caused by pre-commit hook" +
			"try committing manually",
		)
	}
	log.Println("✅ initial commit created successfully")
	return nil
}
