//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

var Default = Build

// Build compiles the project
func Build() error {
	fmt.Println("Building the project...")
	os.Setenv("CGO_ENABLED", "0")
	cmd := exec.Command("go", "build", "-o", "url-shortener")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs the unit tests
func Test() error {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run starts the URL shortener service
func Run() error {
	fmt.Println("Starting the URL shortener service...")
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning up...")
	err := os.Remove("url-shortener")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Println("Cleanup completed.")
	return nil
}

// Generate runs 'go generate' to generate embedded files before Build or Run
func Generate() error {
	fmt.Println("Running 'go generate'...")
	cmd := exec.Command("go", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// BuildWithGenerate is a custom task that runs Generate before Build
func BuildWithGenerate() error {
	err := Generate()
	if err != nil {
		return err
	}
	return Build()
}
