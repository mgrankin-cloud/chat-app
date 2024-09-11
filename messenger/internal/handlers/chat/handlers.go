package user

import (
	"context"
	"errors"

	"log/slog"

	ssov4 "github.com/mgrankin-cloud/messenger/contract/gen/go/chat"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
	"github.com/mgrankin-cloud/messenger/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov4.UnimplementedChatServer
	chat    Chat
	storage *storage.Storage
	log     *slog.Logger
}

type Chat interface {
	GetChat(
		ctx context.Context,
		chatID int64,
	) (chatName string, photo []byte, chatType ssov4.ChatType, status string, err error)
	UpdateChat(
		ctx context.Context,
		chatID int64,
		chatName string,
		photo []byte,
		status string,
	) (success bool, err error)
	CreateChat(
		ctx context.Context,
		name string,
		photo []byte,
		chatType ssov4.ChatType,
	) (success bool, chatID int64, err error)
	DeleteChat(
		ctx context.Context,
		chatID int64,
	) (success bool, err error)
	MuteChat(
		ctx context.Context,
		chatID int64,
		mute bool,
	) (success bool, err error)
	PinChat(
		ctx context.Context,
		chatID int64,
		pin bool,
	) (success bool, err error)
}

func RegisterChatService(gRPCServer *grpc.Server, chat Chat) {
	ssov4.RegisterChatServer(gRPCServer, &serverAPI{chat: chat})
}

func (s *serverAPI) GetChat(
	ctx context.Context,
	request *ssov4.GetChatRequest,
) (*ssov4.GetChatResponse, error) {

	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	chatName, chatType, photo, status, err := s.chat.GetChat(ctx, request.GetChatId())
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get chat: %v", err)
	}

	return &ssov4.GetChatResponse{
		ChatName: chatName,
		Photo: photo,
		ChatType: chatType,
		Status: status,
	}, nil
}

func (s *serverAPI) CreateChat(
	ctx context.Context,
	request *ssov4.CreateChatRequest,
) (*ssov4.CreateChatResponse, error) {
	if request.GetChatName() == "" {
		return nil, status.Error(codes.InvalidArgument, "chat name is required")
	}

	if request.GetChatType() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat type is required")
	}

	success, chatID, err := s.chat.CreateChat(ctx, request.GetChatName(), request.GetPhoto(), request.GetChatType())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	return &ssov4.CreateChatResponse{
		Success: success,
		ChatId:  chatID,
	}, nil
}

func (s *serverAPI) UpdateChat(
	ctx context.Context,
	request *ssov4.UpdateChatRequest,
) (*ssov4.UpdateChatResponse, error) {
	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	if request.GetChatName() == "" {
		return nil, status.Error(codes.InvalidArgument, "chat name is required")
	}

	success, err := s.chat.UpdateChat(ctx, request.GetChatId(), request.GetChatName(), request.GetPhoto(), request.GetStatus())
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update chat: %v", err)
	}

	return &ssov4.UpdateChatResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) DeleteChat(
	ctx context.Context,
	request *ssov4.DeleteChatRequest,
) (*ssov4.DeleteChatResponse, error) {
	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	success, err := s.chat.DeleteChat(ctx, request.GetChatId())
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}

	return &ssov4.DeleteChatResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) MuteChat(
	ctx context.Context,
	request *ssov4.MuteChatRequest,
) (*ssov4.MuteChatResponse, error) {
	chatID := request.GetChatId()
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get user ID from context: %v", err)
	}

	if chatID == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	s.logAction(userID, chatID, "mute")

	success,err := s.storage.updateChatStatus(chatID, "muted"); 
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to mute chat")
	}

	return &ssov4.MuteChatResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) PinChat(
	ctx context.Context,
	request *ssov4.PinChatRequest,
) (*ssov4.PinChatResponse, error) {
	chatID := request.GetChatId()
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get user ID from context: %v", err)
	}

	if chatID == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	s.logAction(userID, chatID, "pin")

	success,err := s.storage.updateChatStatus(chatID, "pinned");
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to pin chat")
	}

	return &ssov4.PinChatResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) logAction(userID, chatID int64, action string) {
	s.log.Info("User action", slog.Int64("userID", userID), slog.Int64("chatID", chatID), slog.String("action", action))
}
