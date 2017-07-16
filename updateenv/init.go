/*

updateenv is a helper package that should be imported for its side effect by a program (main package)
that wants to experiment with the auto-GOPATH semantics proposed by myitcv.io/go

*/
package updateenv

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	FilenameGOPATH  = "._gopath"
	DirectoryVendor = "_vendor"
)

func init() {
	// default value; we need to re-implement the $HOME/go logic here in case any vendors
	// need to be added to the GOPATH
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		// TODO this is probably not as strong as it could be
		gopath = os.Getenv("HOME")
	}

	d, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get current working directory: %v", err))
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

	os.Setenv("GOPATH", gopath)
}
