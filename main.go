package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	DEFAULT_LISTEN_ADDRESS = "0.0.0.0"
	DEFAULT_LISTEN_PORT    = 8080
	DEFAULT_PREFIX         = "/"
)

const header = `<html>
<link rel="stylesheet" href="https://raw.githubusercontent.com/sindresorhus/github-markdown-css/gh-pages/github-markdown.css">
<style>
    .markdown-body {
        min-width: 200px;
        max-width: 790px;
        margin: 0 auto;
        padding: 30px;
    }
</style>
<article class="markdown-body">
`

const footer = `
</article>
</html>`

type config struct {
	root string
}

type context struct {
	cfg *config
}

func (ctx *context) pathForPage(page string) string {
	path := path.Join(ctx.cfg.root, page+".md")
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func (ctx *context) PageHandler(w http.ResponseWriter, r *http.Request) {
	path := ctx.pathForPage(mux.Vars(r)["page"])
	if path == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	input, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(header))
	w.Write(blackfriday.MarkdownCommon(input))
	w.Write([]byte(footer))
}

func setupRouter(r *mux.Router, cfg *config) error {
	ctx := &context{cfg: cfg}
	r.HandleFunc("/{page}", ctx.PageHandler).Methods("GET")
	return nil
}

func main() {
	cfg := &config{root: "/tmp/pages"}

	router := mux.NewRouter()
	if err := setupRouter(router.PathPrefix(DEFAULT_PREFIX).Subrouter(), cfg); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", DEFAULT_LISTEN_ADDRESS, DEFAULT_LISTEN_PORT)
	log.Printf("Starting wiki server on http://%s%s", addr, DEFAULT_PREFIX)
	http.Handle("/", router)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
