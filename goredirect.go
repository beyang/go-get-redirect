package goredirect

import (
	"fmt"
	"github.com/gorilla/pat"
	"html/template"
	"net/http"
	"strings"
)

type RepoNamespace struct {
	VCS        string // currently only tested on git
	Scheme     string
	Hostname   string
	PathPrefix string
}

type Mapping struct {
	Prefix  string // needs to have both a leading and trailing '/'
	DstRepo RepoNamespace
}

func isGoGet(req *http.Request) bool {
	if req.URL.Query().Get("go-get") == "1" {
		return true
	}
	return false
}

func NewGoGetHandler(mappings []Mapping, defaultHandler http.Handler) http.Handler {
	goGetHandler := pat.New()
	for _, mapping := range mappings {
		dst := &mapping.DstRepo
		handlerFunc := repoRedirectFunc(dst.Scheme, dst.Hostname, dst.PathPrefix, dst.VCS)
		goGetHandler.Get(fmt.Sprintf("%s{owner:.+}/{repo:.+}", mapping.Prefix), handlerFunc)
		goGetHandler.Get(fmt.Sprintf("%s{owner:.+}/{repo:.+}/", mapping.Prefix), handlerFunc)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if isGoGet(req) {
			goGetHandler.ServeHTTP(w, req)
		} else {
			if defaultHandler != nil {
				defaultHandler.ServeHTTP(w, req)
			} else {
				http.Error(w, "Not a go-get request and no default handler provided", http.StatusNotFound)
			}
		}
	})
}

// Maps from srchost/src-prefix/:owner/:repo -> dsthost/dst-prefix/:owner/:repo
func repoRedirectFunc(dstScheme string, dstHost string, dstPath string, vcs string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		owner := q.Get(":owner")
		repoName := q.Get(":repo")
		srcHost := strings.Split(req.Host, ":")[0]

		goget := q.Get("go-get")

		if goget == "" {
			panic("trying to serve go-get redirect for non-go-get query")
		}

		t, err := template.New("redirectTemplate").Parse(redirectTemplate)
		if err != nil {
			panic("error parsing template: " + err.Error())
		}

		err = t.Execute(w, templateParams{
			Root:         fmt.Sprintf("%s/%s/%s", srcHost, owner, repoName),
			VCS:          vcs,
			RedirectRoot: fmt.Sprintf("%s://%s%s%s/%s", dstScheme, dstHost, dstPath, owner, repoName),
		})
		if err != nil {
			http.Error(w, "template execution error", http.StatusInternalServerError)
		}
	}
}
