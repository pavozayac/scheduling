package shared

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type TestGrpcServerOption = func(*grpc.Server)

func WithFunc(mutator TestGrpcServerOption) TestGrpcServerOption {
	return mutator
}

func SetupTestGrpcServer(t *testing.T, options ...TestGrpcServerOption) *grpc.ClientConn {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	for _, option := range options {
		option(s)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// Use passthrough resolver, the specific hostname is not important
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return conn
}
