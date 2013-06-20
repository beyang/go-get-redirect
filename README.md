go-get-redirect
===============

Package goredirect enables you to redirect `go get` to map from one set of repository URLs to another.

Caveats
=======
- `go get` assumes all hostnames have at least one . in them, so localhost will not work.  You can
  get around this by adding `127.0.0.1 right.here` to your `/etc/hosts`.
- Currently assumes format of all repo roots is at hostname/base/:owner/:repo
