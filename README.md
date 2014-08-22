Package go-get-redirect
===============

[![xrefs](https://sourcegraph.com/api/repos/github.com/beyang/go-get-redirect/badges/xrefs.png)](https://sourcegraph.com/github.com/beyang/go-get-redirect)
[![funcs](https://sourcegraph.com/api/repos/github.com/beyang/go-get-redirect/badges/funcs.png)](https://sourcegraph.com/github.com/beyang/go-get-redirect)
[![top func](https://sourcegraph.com/api/repos/github.com/beyang/go-get-redirect/badges/top-func.png)](https://sourcegraph.com/github.com/beyang/go-get-redirect)

Package goredirect enables you to redirect `go get` to map from one
set of repository URLs to another. You can use it to create a
standalone server or wrap an existing one.

For a quick demo:

    go install github.com/beyang/go-get-redirect/...
    sudo sh -c 'echo "127.0.0.1 right.here" >> /etc/hosts'
    sudo goredirect`
    go get right.here/beyang/go-get-redirect

Observe that `$GOPATH/src/right.here/beyang/go-get-redirect` now holds
identical contents to `$GOPATH/src/github.com/beyang/go-get-redirect`.

For usage examples, look at `cmd/goredirect/cmd.go` and the more
extensive test cases in `goredirect_test.go`.


Notes
=======

- `go get` assumes all hostnames have at least one `.` in them, so `go
  get localhost/repo/path` will not work. In the dev environment, you
  can get around this by adding a dummy host (e.g., `right.here`) to
  `/etc/hosts`.

- For development purposes, you can pass a custom port number with the
  `--port` flag, but `go get` assumes the server listens on port 80
  (or 443 for HTTPS).

