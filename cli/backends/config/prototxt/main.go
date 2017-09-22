package prototxt

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/tim15/wiz/api/backend"
	"io"
)

func Register() backend.ConfigBackend {
	b := backend.ConfigBackend{
		Name:       "prototxt",
		FileSuffix: "prototxt",
		Parse: func(r io.Reader, pb *proto.Message) error {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			s := buf.String() // Does a complete copy of the bytes in the buffer.
			err := proto.UnmarshalText(s, *pb)
			return err
		},
	}
	return b
}
