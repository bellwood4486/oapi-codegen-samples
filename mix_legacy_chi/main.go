package main

import (
	"fmt"
	"net/http"

	"github.com/bellwood4486/oapi-codegen-samples/mix_legacy_chi/oapi"

	"github.com/go-chi/chi/v5"
)

// BlogAPIImpl implements OpenAPI-based endpoints.
type BlogAPIImpl struct{}

func (*BlogAPIImpl) FindPosts(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("This is new posts api"))
}

func main() {
	r := chi.NewRouter()

	// "/articles" is a legacy endpoint sample that is not OpenAPI based.
	r.Get("/articles", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("This is legacy articles api"))
	})

	// Append OpenAPI-based endpoints(ex."/posts") to the existing router.
	oapi.HandlerFromMux(&BlogAPIImpl{}, r)

	addr := ":8000"
	fmt.Printf("listen %s ...\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
