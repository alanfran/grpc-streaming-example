package example

import (
	"errors"
	"fmt"
	"io"

	example "github.com/alanfran/grpc-streaming-example/pkg/example/v1"
)

// gRPC default maximum message size is 4 MB
// Ciro Costa found 1 Mebibyte to be an optimal message size
// https://ops.tips/blog/sending-files-via-grpc/
const MaxBytesPerChunk = 1024 * 1024

type ExampleServer struct {
}

var _ example.ExampleServer = ExampleServer{}

func (e ExampleServer) CreateBigFile(stream example.Example_CreateBigFileServer) error {
	var contents []byte

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving message from stream: %w", err)
		}

		contents = append(contents, msg.GetBigFileChunk()...)
	}

	// then save the file somewhere

	// and send the unary response to the client
	// Note: because these files will exceed 4 MB in size,
	// we will only be sending the metadata back to the client.
	err := stream.SendAndClose(&example.BigFile{
		Name:      "foo",
		SizeBytes: int64(len(contents)),
	})
	if err != nil {
		return fmt.Errorf("error sending response to client: %w", err)
	}

	return nil
}

func (e ExampleServer) GetBigFile(req *example.GetBigFileRequest, stream example.Example_GetBigFileServer) error {
	// open a big file

	// split it into byte chunks

	// stream chunks to the client

	return errors.New("nyi")
}
