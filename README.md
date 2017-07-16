### `myitcv.io/go`

`myitcv.io/go` is a wrapper around the `go` tool that automatically sets the `GOPATH` env variable based on the process'
current directory (which we refer to as `$PWD` for brevity). See the **Goal** section below for motivation.

The `GOPATH` passed to the `go` tool is calculated as follows:

1. The existence of a `._gopath` file in a directory `d1` sets `GOPATH` to `d1`; this is evaluated for each directory
   `d1` on the path from `$PWD` towards `/`. If no such `d1` exists, set `GOPATH` to the `myitcv.io/go` process'
   environment variable `GOPATH` (if set), else `$HOME/go`
2. The existence of a `_vendor` directory in a directory `d2` causes `d2/_vendor` to be prepended to `GOPATH`. This is
   evaluated for each directory `d2` on the path from `d1` (or `/` if `d1` did not exist) towards `$PWD`.

The algorithm is relatively simple; `myitcv.io/go/updateenv`'s `init` should be relatively understandable therefore.

### Usage

`myitcv.io/go` is designed to be used in place of the `go` tool. Indeed it is implemented under the assumption that
`myitcv.io/go` resolves before the `go` tool in your `PATH`.

Therefore, assuming the `go` tool is already on your `PATH`, do something equivalent to the following

```bash
go get -u myitcv.io/go
mkdir -p "$HOME/bin/myitcv_io_go"
go build -o "$HOME/bin/myitcv_io_go/go" myitcv.io/go
```

and then ensure `$HOME/bin/myitcv_io_go` is at the beginning of your `PATH`:

```bash
# in .bashrc or equivalent
export PATH="$HOME/bin/myitcv_io_go:$PATH"

# ...
```

### Goal of this tool

The goal behind this tool is to explore how a true auto-GOPATH implementation might work within Go proper. But there is
an obvious flaw to the approach of "wrapping" the `go` tool: any program _other_ than the `go` tool that depends on a
package that itself uses the value of `GOPATH` will no longer function correctly. Hence `myitcv.io/go` imports
`myitcv.io/go/updateenv` for its side effect where the process' `GOPATH` environment variable is updated. Other programs
should import this package for its side effects and be recompiled to experiment with this change in behaviour.
