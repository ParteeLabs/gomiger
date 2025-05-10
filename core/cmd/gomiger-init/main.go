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
