FROM golang:1.9-alpine

WORKDIR ${CIRCLE_WORKING_DIRECTORY:-"/go/src/github.com/tim15/wiz/"}
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD sh
