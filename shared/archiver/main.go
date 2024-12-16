package main
/**
import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	ssov1 "github.com/mgrankin-cloud/messenger/contract/gen/go/sso"
	ssov2 "github.com/mgrankin-cloud/messenger/contract/gen/go/user"
	config "github.com/mgrankin-cloud/messenger/internal/config/auth"
	st "github.com/mgrankin-cloud/messenger/internal/handlers/user"
	"github.com/mgrankin-cloud/messenger/pkg/utils/lib/archiver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const configPath = "../config/config_local.yaml"

const grpcHost = "localhost:8081"

type SuiteArchiver struct {
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func main() {
	ctx := context.Background()

	var uid int64 = 1

	userResponse, err := st.GetUser(ctx, &ssov2.GetUserRequest{
		UserId: uid,
	})

	uid = userResponse.GetUserId()

	files := []string{}

	err = archiver.SendArchive(nil, uid, files)
	if err != nil {
		fmt.Printf("Ошибка при отправке архива: %v\n", err)
		return
	}

	fmt.Printf("Архив отправлен пользователю %d\n", uid)
}

func NewArchiver() (context.Context, *SuiteArchiver) {
	cfg := MustLoadByPath(configPath)

	ctx, err := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}

	authClient := ssov1.NewAuthClient(cc)

	return ctx, &SuiteArchiver{
		Cfg:        cfg,
		AuthClient: authClient,
	}
}
**/