package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/go-chi/chi/v5"

	"github.com/bellwood4486/oapi-codegen-samples/validation_chi/oapi"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
)

func DefineUUIDFormat() {
	openapi3.DefineStringFormatCallback("uuid", func(uuidStr string) error {
		_, err := uuid.Parse(uuidStr)
		return err
	})
}

func DefinePostalFormat() {
	openapi3.DefineStringFormat("postal", `^[0-9]{3}-[0-9]{4}$`)
}

func main() {
	swagger, err := oapi.GetSwagger()
	if err != nil {
		fmt.Printf("failed to get swagger spec: %v\n", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	openapi3.DefineIPv4Format()
	openapi3.DefineIPv6Format()
	DefineUUIDFormat()
	DefinePostalFormat()

	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	// Use `OapiRequestValidatorWithOptions` if you want to change validator's behavior.
	//r.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
	//	Options: openapi3filter.Options{
	//		MultiError: true,
	//	},
	//}))

	blogAPI := oapi.NewBlogAPI()
	oapi.HandlerFromMux(blogAPI, r)

	addr := ":8000"
	fmt.Printf("listen %s ...\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("failed to listen and serve: %v\n", err)
		os.Exit(1)
	}
}
