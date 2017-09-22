package backend

import (
	"github.com/golang/protobuf/proto"
)

type ConfigBackend struct {
	Name       string
	FileSuffix string
	Parse      func(io.Reader) proto.Message
}
