/*

go is a wrapper around the go tool that automatically sets the GOPATH env
variable based on the process' current directory.

*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	FilenameGOPATH  = "._gopath"
	DirectoryVendor = "_vendor"
)

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

	// default value; we need to re-implement the $HOME/go logic here in case any vendors
	// need to be added to the GOPATH
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		// TODO this is probably not as strong as it could be
		gopath = os.Getenv("HOME")
	}

	d, err := os.Getwd()
	if err != nil {
		failf("failed to get current working directory: %v", err)
	}

	var vendors []string

	for {
		// check for presence of .gopath file first; doesn't make sense to
		// _vendor in same directory

		gpf, err := os.Stat(filepath.Join(d, FilenameGOPATH))
		if err == nil && gpf.Mode().IsRegular() {
			gopath = d
			break
		}

		vd, err := os.Stat(filepath.Join(d, DirectoryVendor))
		if err == nil && vd.IsDir() {
			vendors = append(vendors, filepath.Join(d, DirectoryVendor))
		}

		nd := filepath.Dir(d)

		if nd == d {
			break
		}

		d = nd
	}

	for i := len(vendors) - 1; i >= 0; i-- {
		gopath = vendors[i] + string(filepath.ListSeparator) + gopath
	}

	env := []string{"GOPATH=" + gopath}

	for _, v := range os.Environ() {
		vs := strings.SplitN(v, "=", 2)

		if vs[0] != "GOPATH" {
			env = append(env, v)
		}
	}

	argv := append([]string{rgp}, os.Args[1:]...)

	err = syscall.Exec(argv[0], argv, env)
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
