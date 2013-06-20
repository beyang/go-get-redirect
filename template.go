package goredirect

var redirectTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>Repository</title>
    <meta name="go-import" content="{{.Root}} {{.VCS}} {{.RedirectRoot}}">
  </head>
  <body>
    Content: {{.Root}} {{.VCS}} {{.RedirectRoot}}
  </body>
</html>`

type templateParams struct {
	Root         string
	VCS          string
	RedirectRoot string
}
