//go:build ignore

// Package main is an initialization tool for gomiger, a code generation utility.
//
// It provides functionality to initialize source code based on a configuration file (gomiger.rc).
// The program checks if the source code has already been initialized, and if not,
// proceeds with the initialization process using the specified configuration.
//
// If any errors occur during the process (such as failing to load the configuration
// file or initialization failures), the program will terminate with an appropriate
// error message.
//
// This tool is typically used as part of the initial setup process for a gomiger-based
// project and should be run only once per project.
package main

import (
	"log"

	"github.com/ParteeLabs/gomiger/core"
	"github.com/ParteeLabs/gomiger/core/generator"
)

var rcPath string

func main() {
	rc, err := core.GetGomigerRC(rcPath)
	if err != nil {
		log.Fatalf("Cannot load the gomiger.rc file: %s", err)
	}
	if generator.IsSrcCodeInitialized(rc) {
		log.Fatalf("The source code is ALREADY INITIALIZED")
	}
	if err := generator.InitSrcCode(rc); err != nil {
		log.Fatalf("Cannot init gomiger: %s", err)
	}
}
