package util

import (
	"os"
	"path/filepath"
	"fmt"
	"os/exec"
)

const (
	goBuildFile       = "build.go"
	fileShimGo        = "shim.go"
	fileShimSupportGo = "shim_support.go"
	dirShim           = "shim"
)

func CreateShim(dependency *Dependency, shimType string) error {

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	createShimSupportGoFile(pwd)

	//todo blow up if shim file not found
	shimFilePath := filepath.Join(dependency.Dir, dirShim, fileShimGo)
	fmt.Println("Shim File:", shimFilePath)

	CopyFile(shimFilePath, filepath.Join(pwd, fileShimGo))

	// Check if this shim based trigger has a gobuild file.
	goBuildFilePath := filepath.Join(dependency.Dir, dirShim, goBuildFile)
	if _, err := os.Stat(goBuildFilePath); err == nil {
		fmt.Println("This trigger makes use if a gobuild file...")
		fmt.Println("Go build file:", goBuildFilePath)
		//fgutil.CopyFile(gobuildFilePath, filepath.Join(env.GetAppDir(), gobuildFile))

		// Execute go run gobuild.go
		cmd := exec.Command("go", "run", goBuildFilePath, pwd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		//cmd.Dir = env.GetAppDir()
		//cmd.Env = append(os.Environ(),
		//	fmt.Sprintf("GOPATH=%s", env.GetRootDir()),
		//)

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func createShimSupportGoFile(codeSourcePath string) {

	f, _ := os.Create(filepath.Join(codeSourcePath, fileShimSupportGo))
	RenderTemplate(f, tplShimSupportGoFile, nil)
	f.Close()
}

//func removeShimGoFiles(codeSourcePath string) {
//	os.Remove(filepath.Join(codeSourcePath, fileShimGo))
//	os.Remove(filepath.Join(codeSourcePath, fileShimSupportGo))
//}

var tplShimSupportGoFile = `// Do not change this file, it has been generated
// If you change it and rebuild the application your changes might get lost
package main

import (
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
)

func init() {
	config.SetDefaultLogLevel("ERROR")
	logger.SetLogLevel(logger.ErrorLevel)

	app := shimApp()

	e, err := flogo.NewEngine(app)

	if err != nil {
		logger.Error(err)
		return
	}

	e.Init(true)
}
`
