package main

import (
	"flag"
	"fmt"
	goredirect "github.com/beyang/go-get-redirect"
	"log"
	"net/http"
)

var port = flag.Int("port", 80, "HTTP listen port.  You can set this to something other than 80 for debugging purposes, but it has to be 80 when using go get.")

func main() {
	flag.Parse()

	http.Handle("/", goredirect.NewGoGetHandler([]goredirect.Mapping{
		{"git", "https", "github.com", goredirect.NewStringMapperOrBust("/(?P<owner>.+)/(?P<repo>.+)", "/{{.owner}}/{{.repo}}")},
	}, nil))

	log.Printf("Starting server on 0.0.0.0:%d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
