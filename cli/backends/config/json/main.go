package json

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/tim15/wiz/api/backend"
	"io"
)

func Register() backend.ConfigBackend {
	fmt.Println("Reg")
	b := backend.ConfigBackend{
		Name:       "json",
		FileSuffix: "json",
		Parse: func(r io.Reader, pb *proto.Message) error {
			err := jsonpb.Unmarshal(r, *pb)
			return err
		},
	}
	return b
}
