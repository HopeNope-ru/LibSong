.PHONY: build
go-build:
	go build

.PHONY: swagger-generate
swagger-gen:
	swagger generate spec -o swagger.yml

.PHONY: swagger-serve
swagger-serve:
	swagger serve --no-open --port 8080 swagger.yml