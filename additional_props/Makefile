.PHONY: oapigen
API_GEN_FILE := "oapi/api.gen.go"

oapigen:
	rm -f $(API_GEN_FILE)
	oapi-codegen -generate types,chi-server -package oapi blog-api.yaml > $(API_GEN_FILE)

