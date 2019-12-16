#!/usr/bin/env bash

bash $(dirname "$0")/gen-swagger.sh

file=swagger.json

echo "Serving Swagger/OpenAPI 2.0 Spec from '$file'"
swagger serve $file