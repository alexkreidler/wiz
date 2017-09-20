#!/usr/bin/env sh
: ${PROJECT_PATH="$GOPATH/src/github.com/tim15/wiz"}
: ${CONTAINER_PROJECT_PATH="/go/src/github.com/tim15/wiz"}
: ${CONTAINER_GOPATH="/go"}
: ${CONTAINER_NAME="wizdev"}

start_container() {
  # Start dev/builder container
  if [ ! "$(docker ps -a | grep "$CONTAINER_NAME")" ]; then
    cd "$PROJECT_PATH"
    docker create -it -v $GOPATH:"$CONTAINER_GOPATH" --name $CONTAINER_NAME tim15/goproto
  fi;

  docker start $CONTAINER_NAME > /dev/null
}

dev() {
  docker exec -it $CONTAINER_NAME sh
}

genproto() {
  docker exec "$CONTAINER_NAME" $CONTAINER_PROJECT_PATH/scripts/genproto.sh $CONTAINER_PROJECT_PATH
}


# Major issue: Syncing GOPATH will lead to incompatible compiled binaries for different OSes
# Fix: just compile on host?
# Or envourage all dev in container?
# build() {
#   echo "Building + installing binaries"
#   docker exec $CONTAINER_NAME go install ./...
# }
#
# wb() {
#   echo "Watch and rebuild binaries on changes ..."
#   docker exec $CONTAINER_NAME scripts/wb.sh
# }
#
# stop() {
#   echo "Stopping wiz project container(s) ..."
#   docker stop -t 1 $CONTAINER_NAME
#   docker rm $CONTAINER_NAME
# }

usage() {
  echo "
Usage: tool.sh [COMMAND]
Commands:
  dev - start a development container
  genproto - generate protobuf code from the api
  build - build and install the binaries
  wb - watch and build on changes
  stop - stop and remove containers used by this tool
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
    build) start_container; build;;
    wb) start_container; wb;;
    stop) stop;;
    *) usage
  esac
done
