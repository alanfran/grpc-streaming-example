package example

import (
	"errors"

	example "github.com/alanfran/grpc-streaming-example/pkg/example/v1"
)

type ExampleServer struct {
}

var _ example.ExampleServer = ExampleServer{}

func (e ExampleServer) CreateBigFile(stream example.Example_CreateBigFileServer) error {
	return errors.New("nyi")
}

func (e ExampleServer) GetBigFile(req *example.GetBigFileRequest, stream example.Example_GetBigFileServer) error {
	return errors.New("nyi")
}
