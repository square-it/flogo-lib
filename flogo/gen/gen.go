package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-lib/flogo/gen/util"
)

func main() {
	shim := flag.Bool("shim", false, "generate shim")
	flag.Parse()
	args := os.Args

	fmt.Println(os.Getwd())

	path, _ := os.Getwd()

	if len(args) > 1 {

		if args[1] == "-shim" {
			if len(args) > 2 {
				path = args[2]
			}
		} else {
			path = args[1]
		}
	}

	if *shim {
		generateShim(path)
		return
	}

	util.GenerateGoMetadata(args[1])
}

func generateShim(path string) {

	dependencies, err := util.GenerateGoMetadata(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for _, dependency := range dependencies {
		if dependency.ContribType == util.TRIGGER {
			mdFilePath := filepath.Join(dependency.Dir, "trigger.json")

			mdBytes, err := ioutil.ReadFile(mdFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			md, err := util.ParseTriggerMetadata(mdBytes)

			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			if md.Shim != "" {
				util.CreateShim(dependency, md.Shim)
			}
		}
	}
}
