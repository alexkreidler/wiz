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
