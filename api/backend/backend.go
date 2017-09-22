package backend

import (
	"github.com/golang/protobuf/proto"
	"io"
)

type ConfigBackend struct {
	Name       string
	FileSuffix string
	Parse      func(r io.Reader, pb *proto.Message) error
}
