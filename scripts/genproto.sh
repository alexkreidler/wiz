#!/usr/bin/env sh
if [ "$1" != "" ]; then
  : ${PROJECT_PATH:="$1"}
else
  : ${PROJECT_PATH:="$(pwd)"}
fi

if grep docker /proc/1/cgroup -qa; then
    echo "Compiling protobufs"
else
    echo "WARNING: Not in container. It is recommended to use the tim15/goproto image for development"
fi

set -o xtrace
# protoc -I "$GOPATH/src" --go_out=plugins=grpc:$GOPATH/src $rawSpec $PROJECT_PATH/api/proto/*.proto

for rawSpec in $PROJECT_PATH/api/proto/*.proto; do
  # uses `option go_package` to specify absolute location from $GOPATH/bin
  protoc -I "$GOPATH/src" --go_out=plugins=grpc:$GOPATH/src $rawSpec
done

echo "All done!"
