echo $1
PROJECT_PATH="$1"

if grep docker /proc/1/cgroup -qa; then
    echo "Compiling protobufs"
else
    echo "Not in container. It is recommended to use the tim15/goproto image for development"
fi

protoc -I "$PROJECT_PATH/api/proto" --go_out="$PROJECT_PATH/api/gen" $PROJECT_PATH/api/proto/*
