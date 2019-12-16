#!/usr/bin/env bash

file=processors.apib

echo "Generating Swagger/OpenAPI 2.0 Spec from '$file'"
apib2swagger -i $file -o swagger.json