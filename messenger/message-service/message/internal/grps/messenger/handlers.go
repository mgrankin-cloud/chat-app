package user

import (
	ssov4 "Messenger-android/messenger/chat-service/contract/gen/go/chat"
	"Messenger-android/messenger/internal/storage"
	ssov5 "Messenger-android/messenger/message-service/contract/gen/go/message"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov5.UnimplementedMessageServer
	message Message
	storage *storage.Storage
}

type Message interface {
	GetMessage(
		ctx context.Context,
		request *ssov5.GetMessageRequest,
	) (*ssov4.GetChatResponse, error)
	UpdateMessage(
		ctx context.Context,
		request *ssov5.UpdateMessageRequest,
	) (*ssov5.UpdateMessageResponse, error)
	DeleteMessage(
		ctx context.Context,
		request *ssov5.DeleteMessageRequest,
	) (*ssov5.DeleteMessageResponse, error)
	PinMessage(
		ctx context.Context,
		request *ssov5.PinMessageRequest,
	) (*ssov5.PinMessageResponse, error)
}

func RegisterMessageService(gRPCServer *grpc.Server, message Message) {
	ssov5.RegisterMessageServer(gRPCServer, &serverAPI{message: message})
}

func (s *serverAPI) GetMessage(
	ctx context.Context,
	request *ssov5.GetMessageRequest,
) (*ssov5.GetMessageResponse, error) {

	if request.GetMessageId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	message, err := s.storage.GetMessageByID(ctx, request.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	return &ssov5.GetMessageResponse{}, nil
}

func (s *serverAPI) UpdateMessage(
	ctx context.Context,
	request *ssov5.UpdateMessageRequest,
) (*ssov5.UpdateMessageResponse, error) {
	messageID := request.GetMessageId()

	if messageID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	err := s.storage.UpdateMessage(ctx, messageID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update message: %v", err)
	}

	return &ssov5.UpdateMessageResponse{}, nil
}

func (s *serverAPI) DeleteMessage(
	ctx context.Context,
	request *ssov5.DeleteMessageRequest,
) (*ssov5.DeleteMessageResponse, error) {
	messageID := request.GetMessageId()

	if messageID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	err := s.storage.DeleteMessage(ctx, messageID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete message: %v", err)
	}

	return &ssov5.DeleteMessageResponse{}, nil
}

func (s *serverAPI) PinMessage(
	ctx context.Context,
	request *ssov5.PinMessageRequest,
) (*ssov5.PinMessageResponse, error) {
	messageID := request.GetMessageId()
	userID := storage.getUserIDFromContext(ctx)

	if messageID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	s.lib.logAction(userID, messageID, "pin")

	if err := s.lib.updateChatStatus(messageID, "pinned"); err != nil {
		return nil, status.Error(codes.Internal, "failed to pin message")
	}

	s.lib.notifyChatMembers(messageID, "Message has been pinned")

	return &ssov5.PinMessageResponse{
		Success: true,
	}, nil
}
