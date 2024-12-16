package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/mgrankin-cloud/messenger/contract/gen/go/sso"
	"github.com/mgrankin-cloud/messenger/internal/config/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const configPath = "../../config/config_auth.yaml"

const grpcHost = "localhost:8081"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath(configPath)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	authClient := ssov1.NewAuthClient(cc)

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authClient,
	}
}
