#!/usr/bin/env sh
: ${PROJECT_PATH="$GOPATH/src/github.com/tim15/wiz"}
: ${CONTAINER_PROJECT_PATH="/go/src/github.com/tim15/wiz"}
: ${CONTAINER_NAME="wizdev"}

function start_container() {
  # Start dev/builder container
  if [ ! "$(docker ps -a | grep "$CONTAINER_NAME")" ]; then
    cd "$PROJECT_PATH"
    docker create -it -v $(pwd):"$CONTAINER_PROJECT_PATH" --name $CONTAINER_NAME tim15/goproto
  fi;

  docker start $CONTAINER_NAME > /dev/null
}

function dev() {
  docker exec -it $CONTAINER_NAME sh
}

function genproto() {
  docker exec "$CONTAINER_NAME" $CONTAINER_PROJECT_PATH/scripts/genproto.sh $CONTAINER_PROJECT_PATH
}

function usage() {
  echo "
Usage: tool.sh [COMMAND]
Commands:
  dev - start a development container
  genproto - generate protobuf code from the api
  "
  exit 1
}

if [[ "$@" == "" ]]; then
  usage
fi

for i in $@; do
  case "$i" in
    dev) echo "Wiz Project dev container"; start_container; dev;;
    genproto) start_container; genproto;;
    *) usage
  esac
done
