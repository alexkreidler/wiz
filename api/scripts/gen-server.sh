#!/usr/bin/env bash

./gen-swagger

path=../../../

echo "cleaning old code"
rm -rf $path/models
rm -rf $path/restapi

echo "Starting server code generation"
swagger generate server -t $path