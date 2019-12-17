API_DIR=./api
BLUEPRINT_FILE=$(API_DIR)/processors.apib
SWAGGER_FILE=$(API_DIR)/swagger.json

start-blueprint:
	@echo Starting "aglio" API Blueprint server
	aglio -i $(BLUEPRINT_FILE) -s

gen-swagger:
	@echo Generating Swagger/OpenAPI 2.0 Spec
	apib2swagger -i $(BLUEPRINT_FILE) -o $(SWAGGER_FILE)

start-swagger: gen-swagger
	@echo Serving Swagger/OpenAPI 2.0 Spec
	swagger serve $(SWAGGER_FILE)

clean-swagger:
	rm -rf models
	rm -rf restapi

clean-gen-swagger-server: clean-swagger gen-swagger-server

gen-swagger-server: # in the past we did this, now swagger is the truth: gen-swagger
	swagger generate server -f $(SWAGGER_FILE)



generate:
	# Generate gogo, gRPC-Gateway, swagger, go-validators output.
	#
	# -I declares import folders, in order of importance
	# This is how proto resolves the protofile imports.
	# It will check for the protofile relative to each of these
	# folders and use the first one it finds.
	#
	# --gogo_out generates GoGo Protobuf output with gRPC plugin enabled.
	# --grpc-gateway_out generates gRPC-Gateway output.
	# --swagger_out generates an OpenAPI 2.0 specification for our gRPC-Gateway endpoints.
	# --govalidators_out generates Go validation files for our messages types, if specified.
	#
	# The lines starting with Mgoogle/... are proto import replacements,
	# which cause the generated file to import the specified packages
	# instead of the go_package's declared by the imported protof files.
	#
	# $$GOPATH/src is the output directory. It is relative to the GOPATH/src directory
	# since we've specified a go_package option relative to that directory.

	# We use gogoslick for development and printing nice structs, but we can move to gogofaster to get smaller output sizes.
	protoc \
		-I proto \
		-I $(GOPATH)/src \
		--gogoslick_out=plugins=grpc,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types:./api ./proto/*.proto

install:
	go get \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/gogo/protobuf/protoc-gen-gogoslick \
		github.com/gogo/protobuf/protoc-gen-gogofaster