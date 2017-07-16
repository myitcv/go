/*

go is a wrapper around the go tool that automatically sets the GOPATH env
variable based on the process' current directory.

*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	_ "myitcv.io/go/updateenv"
)

var debug = false

func main() {
	// find the _second_ uniq go executable in PATH; the first is assumed to be
	// this package (dgp = dummy go path). The second (rgp = real go path) is
	// assumed to be the path to the go exec
	//
	// TODO might need to improve this logic? For example it's easy to get into
	// a loop if myitcv.io/react resolves in the first and second place in
	// PATH... switch instead to logic that requires something that answers `go
	// whyenv` first... followed and then finds the next go executable in PATH
	// that does _not_ answer `go whyenv`

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-debug" {
		debug = true
		args = args[1:]
	}

	pathList := filepath.SplitList(os.Getenv("PATH"))

	dgp, pathList, err := lookPath("go", pathList)
	if err != nil {
		failf("failed to find dummy go in PATH\n")
	}

	var rgp string

	for {
		gp, npathList, err := lookPath("go", pathList)

		if err != nil {
			failf("failed to find real go in PATH remainder %v\n", filepath.Join(pathList...))
		}

		if gp != dgp {
			rgp = gp
			break
		}

		pathList = npathList
	}

	argv := append([]string{rgp}, args...)

	debugf("%v: %v (%v)\n", argv[0], argv[1:], os.Getenv("GOPATH"))

	err = syscall.Exec(argv[0], argv, os.Environ())
	if err != nil {
		failf("failed exec %v: %v\n", rgp, err)
	}
}

func lookPath(file string, pathList []string) (string, []string, error) {
	for i, dir := range pathList {
		if dir == "" {
			dir = "."
		}
		path := filepath.Join(dir, file)
		if err := findExecutable(path); err == nil {
			return path, pathList[i+1:], nil
		}
	}

	return "", nil, os.ErrNotExist
}

func failf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func debugf(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
		os.Stdout.Sync()
	}
}
