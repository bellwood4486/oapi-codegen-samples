.PHONY: oapigen
API_GEN_FILE := "api/api.gen.go"

oapigen:
	rm -f $(API_GEN_FILE)
	oapi-codegen -generate types,chi-server blog-api.yaml > $(API_GEN_FILE)
