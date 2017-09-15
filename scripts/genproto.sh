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

protoc -I "$PROJECT_PATH/api/proto" --go_out="$PROJECT_PATH/api/gen" $PROJECT_PATH/api/proto/*
echo "Done"
