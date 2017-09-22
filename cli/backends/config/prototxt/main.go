package prototxt

import (
  "github.com/golang/protobuf/proto"
  "github.com/wiz/cli/backends"
  "io"
  "bytes"
)

func Register() (backends.ConfigBackend) {
  b := &backends.ConfigBackend{
    Name: "prototxt",
    FileSuffix: "prototxt",
    Parse: func(r io.Reader, pb *proto.Message) error {
      buf := new(bytes.Buffer)
      buf.ReadFrom(r)
      s := buf.String() // Does a complete copy of the bytes in the buffer.
      err := proto.UnmarshalText(s, pb)
      return err
    }
  }
}
