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

.PHONY: proto
proto:
	mkdir ./api
	protoc --proto_path=./proto --go_out=./api ./proto/*.proto

clean-proto:
	rm -rf ./api/