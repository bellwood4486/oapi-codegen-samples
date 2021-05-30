package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/go-chi/chi/v5"

	"github.com/bellwood4486/oapi-codegen-samples/validation_chi/oapi"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
)

func main() {
	swagger, err := oapi.GetSwagger()
	if err != nil {
		fmt.Printf("failed to get swagger spec: %v\n", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	openapi3.DefineIPv4Format()
	openapi3.DefineIPv6Format()

	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))

	blogAPI := oapi.NewBlogAPI()
	oapi.HandlerFromMux(blogAPI, r)

	addr := ":8000"
	fmt.Printf("listen %s ...\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("failed to listen and serve: %v\n", err)
		os.Exit(1)
	}
}
