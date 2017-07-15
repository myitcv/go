### `myitcv.io/go`

`myitcv.io/go` is a wrapper around the `go` tool that automatically sets the `GOPATH` env variable based on the process'
current directory (which we refer to as `$PWD` for brevity).

The `GOPATH` passed to the `go` tool is calculated as follows:

1. The existence of a `._gopath` file in a directory `d1` sets `GOPATH` to `d1`; this is evaluated for each directory
   `d1` on the path from `$PWD` towards `/`. If no such `d1` exists, set `GOPATH` to the `myitcv.io/go` process'
   environment variable `GOPATH` (if set), else `$HOME/go`
2. The existence of a `_vendor` directory in a directory `d2` causes `d2/_vendor` to be prepended to `GOPATH`. This is
   evaluated for each directory `d2` on the path from `d1` (or `/` if `d1` did not exist) towards `$PWD`.

The algorithm is relatively simple; `main.go` should be relatively understandable therefore.

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
