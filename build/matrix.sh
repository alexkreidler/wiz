declare oses=(darwin windows) arches=(amd64 386)
for os in ${oses[*]}; do
  for arch in ${arches[*]}; do
    echo $os $arch
    ./build.sh $os $arch
  done
done

echo hi
os=""
declare arches=(amd64 arm arm64)
for arch in ${arches[*]}; do
  echo $arch
  if [[ "$arch" == "amd64" ]]; then
    ./build.sh "linux" "amd64" "latest"
  else
    ./build.sh "linux" $arch
  fi
done
