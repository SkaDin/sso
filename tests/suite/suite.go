package suite

import (
	"context"
	ssov1 "github.com/SkaDin/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"sso/internal/config"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})
	cc, err := grpc.NewClient(grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc server connection")
	}
	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"
	if v := os.Getenv(key); v != "" {
		return v
	}
	return "../config/local.yml"
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.Grpc.Port))
}
