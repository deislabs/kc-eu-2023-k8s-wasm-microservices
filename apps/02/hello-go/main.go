package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/julienschmidt/httprouter"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		router := spinhttp.NewRouter()
		router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "Hello, WebAssembly!")
		})
		router.ServeHTTP(w, r)
	})
}

func main() {}
