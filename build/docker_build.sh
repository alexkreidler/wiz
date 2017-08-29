declare GOOS=$1 GOARCH=$2 TAG=$3
docker build ./out/$GOOS-$GOARCH -t wizproject/wiz:${TAG:-"$GOOS-$GOARCH"}
if [[ "$CI" == "true" ]]; then
  docker login -u $DOCKER_USER -p $DOCKER_PASS
  docker push wizproject/wiz:${TAG:-"$GOOS-$GOARCH"}
fi
