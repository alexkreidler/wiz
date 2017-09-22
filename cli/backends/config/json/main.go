package json

import (
  "github.com/golang/protobuf/jsonpb"
  "github.com/golang/protobuf/proto"
  "github.com/wiz/cli/backends"
  "io"
)

func Register() (backends.ConfigBackend) {
  b := &backends.ConfigBackend{
    Name: "json",
    FileSuffix: "json",
    Parse: func(r io.Reader, pb *proto.Message) error {
      err := jsonpb.Unmarshal(r, pb)
      return err
    }
  }
}
