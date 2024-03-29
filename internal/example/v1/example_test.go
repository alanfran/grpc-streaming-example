package example_test

import (
	"context"
	"crypto/rand"
	"io"
	"testing"

	example "github.com/alanfran/grpc-streaming-example/internal/example/v1"
	exampleproto "github.com/alanfran/grpc-streaming-example/pkg/example/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockCreateBigFileStream mocks the stream that CreateBigFile will consume.
type MockCreateBigFileStream struct {
	// embed this interface to make the type checker happy
	grpc.ServerStream

	// input bytes split into chunks for streaming
	chunks [][]byte

	// keep track of progress when consuming the stream
	i *int

	// metadata returned from CreateBigFile
	response *exampleproto.BigFile
}

func (MockCreateBigFileStream) Context() context.Context {
	return context.Background()
}

// Recv returns the next message from the queue, or io.EOF if the stream is finished.
func (m MockCreateBigFileStream) Recv() (*exampleproto.CreateBigFileRequest, error) {
	if *m.i == len(m.chunks) {
		return nil, io.EOF
	}

	response := &exampleproto.CreateBigFileRequest{
		BigFileChunk: m.chunks[*m.i],
	}

	*m.i++

	return response, nil
}

// SendAndClose is the method the server uses to send the unary response to the client
// and close the stream.
func (m MockCreateBigFileStream) SendAndClose(msg *exampleproto.BigFile) error {
	*m.response = *msg
	return nil
}

func NewCreateBigFileStream(contents []byte) MockCreateBigFileStream {
	// iterator used internally to track the stream progress
	i := 0

	// split up the contents into chunks
	// like a client would do
	var chunks [][]byte
	var chunk []byte

	for i := range contents {
		chunk = append(chunk, contents[i])

		if len(chunk) >= example.MaxBytesPerChunk {
			chunks = append(chunks, chunk)
			chunk = []byte{}
		}
	}

	// save the last chunk
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}

	return MockCreateBigFileStream{
		i:        &i,
		chunks:   chunks,
		response: &exampleproto.BigFile{},
	}
}

func TestCreateBigFile(t *testing.T) {
	cases := map[string]struct {
		nBytes int64 // number of randomly-generated bytes to write
	}{
		"empty file": {
			nBytes: 0,
		},
		"4 KiB file": {
			nBytes: 4 * 1024,
		},
		"4 MiB file": {
			nBytes: 4 * 1024 * 1024,
		},
		"12 MiB file": {
			nBytes: 12 * 1024 * 1024,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// generate random bytes
			randomBytes := make([]byte, c.nBytes)

			n, err := rand.Read(randomBytes)
			require.Nil(t, err)
			require.Equal(t, c.nBytes, int64(n))

			// create server
			server := example.ExampleServer{}

			// create the mocked stream
			stream := NewCreateBigFileStream(randomBytes)

			// call CreateBigFile with the mocked stream
			err = server.CreateBigFile(stream)
			require.Nil(t, err)

			// assert the returned metadata has the right file size
			require.Equal(t, c.nBytes, stream.response.GetSizeBytes())
		})
	}
}

func TestGetBigFile(t *testing.T) {

}
