package main

import (
	"fmt"
	"github.com/bellwood4486/oapi-codegen-samples/additional_props/oapi"
	"net/http"
)

func main() {
	r := oapi.Handler(oapi.NewBlogAPIImpl())
	addr := ":8000"
	fmt.Printf("listen %s ...\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
