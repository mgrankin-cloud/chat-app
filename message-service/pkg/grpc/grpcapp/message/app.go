package grpcapp

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"strings"
	msggrpc "github.com/mgrankin-cloud/messenger/message-service/internal/handlers/message"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int // Порт, на котором будет работать grpc-сервер
}

// adapts slog logger to interceptor logger.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, fields...)
		case logging.LevelInfo:
			l.Info(msg, fields...)
		case logging.LevelWarn:
			l.Warn(msg, fields...)
		case logging.LevelError:
			l.Error(msg, fields...)
		default:
			l.Info(msg, fields...)
		}
	})
}

// masks sensitive fields in the given message.
func MaskSensitiveFields(msg proto.Message) proto.Message {
	if msg == nil {
		return msg
	}

	// Get the reflection interface for the message
	m := msg.ProtoReflect()

	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if strings.ToLower(string(fd.Name())) == "password" {
			m.Set(fd, protoreflect.ValueOfString("*****"))
		}
		return true
	})

	return msg
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	messageService msggrpc.Message,
	port int,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	msggrpc.RegisterMessageService(gRPCServer, messageService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) RegisterService(registerFunc func(*grpc.Server)) {
	registerFunc(a.gRPCServer)
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	// создание listener, который будет слушить TCP-сообщения
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("sso grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() error {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping sso gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

	return nil
}
