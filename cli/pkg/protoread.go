package pkg

import (
	proto "github.com/golang/protobuf/proto"
	pb "github.com/tim15/wiz/api/proto"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//ProtoRead reads a proto file and returns the pb.Package structure
func ProtoRead(file string) *pb.Package {
	buf, err := ioutil.ReadFile(file)
	check(err)
	dat := &pb.Package{}
	err = proto.Unmarshal(buf, dat)
	check(err)
	return dat
}

//ProtoWrite writes a proto file
func ProtoWrite(file string, pb *pb.Package) error {
	dat, err := proto.Marshal(pb)
	check(err)
	err = ioutil.WriteFile(file, dat, os.ModePerm)
	return err
}
