SWAGGER_FILE=./docs/swagger.json

gen-swagger:
	swag init -g ./server/server.go -p pascalcase

start-swagger: gen-swagger
	@echo Serving Swagger/OpenAPI 2.0 Spec
	swagger serve $(SWAGGER_FILE)

gen-swagger-client: gen-swagger
	swagger generate client -f docs/swagger.json

clean-swagger:
	rm -rf models
	rm -rf client
	rm -rf restapi

add-proto-binaries:
	# Kubernetes go to protobuf generator needs the following binaries:
	# go get golang.org/x/tools/cmd/goimports
	go get github.com/gogo/protobuf/proto
	# can add protoc-gen-gogofast
	go get github.com/gogo/protobuf/protoc-gen-gogo
	go get github.com/gogo/protobuf/gogoproto
	# TODO: look at proteus and kubernetes/code-generator for generate proto files from go APIs
	# would prefer the kubernetes-style method where one HTTP/2 server serves both HTTP and gRPC/proto endpoints with different mime types
