# usage: ./build.sh OS ARCH [TAG]

declare GOOS=$1 GOARCH=$2 TAG=$3
echo "Building wiz for OS: $GOOS on ARCH: $GOARCH"
go build -o out/$GOOS-$GOARCH/wiz ../...
cat << EOF > out/$GOOS-$GOARCH/Dockerfile
FROM scratch

COPY wiz /

CMD ["/wiz"]
EOF
