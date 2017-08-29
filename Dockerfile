FROM golang:1.9-alpine

WORKDIR /go/src/wiz
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["bash"]
