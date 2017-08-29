if [[ "$BUILD_TYPE" == "docker" ]]; then
  script="./docker_build.sh"
else
  script="./build.sh"
fi

declare oses=(darwin windows) arches=(amd64 386)
for os in ${oses[*]}; do
  for arch in ${arches[*]}; do
    echo $os $arch
    $script $os $arch
  done
done

echo hi
os=""
declare arches=(amd64 arm arm64)
for arch in ${arches[*]}; do
  echo $arch
  if [[ "$arch" == "amd64" ]]; then
    $script "linux" "amd64" "latest"
  else
    $script "linux" $arch
  fi
done
