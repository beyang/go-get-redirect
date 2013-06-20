package goredirect

var RedirectTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>Repository</title>
    <meta name="go-import" content="{{.Root}} {{.VCS}} {{.RedirectRoot}}">
  </head>
  <body>
    Content: {{.Root}} {{.VCS}} {{.RedirectRoot}}
  </body>
</html>`
