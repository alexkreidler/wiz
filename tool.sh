#!/usr/bin/env sh
: ${PROJECT_PATH="$GOPATH/src/github.com/tim15/wiz"}
: ${CONTAINER_PROJECT_PATH="/go/src/github.com/tim15/wiz"}
: ${CONTAINER_NAME="wizdev"}

start_container() {
  # Start dev/builder container
  if [ ! "$(docker ps -a | grep "$CONTAINER_NAME")" ]; then
    cd "$PROJECT_PATH"
    docker create -it -v $(pwd):"$CONTAINER_PROJECT_PATH" --name $CONTAINER_NAME tim15/goproto
  fi;

  docker start $CONTAINER_NAME > /dev/null
}

dev() {
  docker exec -it $CONTAINER_NAME sh
}

genproto() {
  docker exec "$CONTAINER_NAME" $CONTAINER_PROJECT_PATH/scripts/genproto.sh $CONTAINER_PROJECT_PATH
}

watch() {
  ./scripts/build_watcher.sh
}

stop() {
  docker stop -t 1 "$CONTAINER_NAME"
}

usage() {
  echo "
Usage: tool.sh [COMMAND]
Commands:
  dev - start a development container
  genproto - generate protobuf code from the api
  watch - watch for changes and rebuild code
  stop - stop the container used by this tool
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
    watch) watch;;
    stop) stop;;
    *) usage
  esac
done
