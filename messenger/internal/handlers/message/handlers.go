package user

import (
	"context"
	"errors"
	"fmt"

	ssov5 "github.com/mgrankin-cloud/messenger/contract/gen/go/message"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
	"github.com/mgrankin-cloud/messenger/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverAPI struct {
	ssov5.UnimplementedMessageServer
	message Message
	storage *storage.Storage
}

type Message interface {
	GetMessage(
		ctx context.Context,
		messageID int64,
	) (content string, photo []byte, createdAt timestamppb.Timestamp, createdBy int64, replyToID int64, receivedAt timestamppb.Timestamp, receivedBy int64, status string, isRead bool, err error)
	CreateMessage(
		ctx context.Context,
		content string,
		createdAt *timestamppb.Timestamp,
		createdBy int64,
		replyToID int64,
		receivedAt int64,
		receivedBy *timestamppb.Timestamp,
	) (success bool, messageID int64, err error)
	UpdateMessage(
		ctx context.Context,
		messageID int64,
		content string,
		status string,
		isRead bool,
	) (success bool, err error)
	DeleteMessage(
		ctx context.Context,
		messageID int64,
	) (success bool, err error)
	PinMessage(
		ctx context.Context,
		messageID int64,
		pin bool,
	) (success bool, err error)
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

	_, err := s.storage.GetMessageByID(ctx, request.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	return &ssov5.GetMessageResponse{
		MessageId: request.GetMessageId(),
	}, nil
}

func (s *serverAPI) CreateMessage(
	ctx context.Context,
	request *ssov5.CreateMessageRequest,
) (*ssov5.CreateMessageResponse, error) {

	if request.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "message content is required")
	}

	success, messageID, err := s.message.CreateMessage(ctx, request.GetContent(), request.GetCreatedAt(), request.GetCreatedBy(), request.GetReplyTo(), request.GetReceivedBy(), request.GetReceivedAt())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create message")
	}

	return &ssov5.CreateMessageResponse{
		Success:   success,
		MessageId: messageID,
	}, nil
}

func (s *serverAPI) UpdateMessage(
	ctx context.Context,
	request *ssov5.UpdateMessageRequest,
) (*ssov5.UpdateMessageResponse, error) {

	if request.GetMessageId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	if request.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "message content is required")
	}

	success, err := s.message.UpdateMessage(ctx, request.GetMessageId(), request.GetContent(), request.GetStatus(), request.GetIsRead())
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update message: %v", err)
	}

	return &ssov5.UpdateMessageResponse{
		Success:   success,
	}, nil
}

func (s *serverAPI) DeleteMessage(
	ctx context.Context,
	request *ssov5.DeleteMessageRequest,
) (*ssov5.DeleteMessageResponse, error) {
	if request.GetMessageId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	success, err := s.message.DeleteMessage(ctx, request.GetMessageId())
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete message: %v", err)
	}

	return &ssov5.DeleteMessageResponse{
		Success:   success,
	}, nil
}

func (s *serverAPI) PinMessage(
	ctx context.Context,
	request *ssov5.PinMessageRequest,
) (*ssov5.PinMessageResponse, error) {
	messageID := request.GetMessageId()
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user id from context: %v", err)
	}

	if messageID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	s.logAction(userID, messageID, "pin")

	if err := s.storage.updateMessageStatus(messageID, "pinned"); err != nil {
		return nil, status.Error(codes.Internal, "failed to pin message")
	}

	return &ssov5.PinMessageResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) logAction(userID, chatID int64, action string) {
	fmt.Printf("User %d performed %s on message %d\n", userID, action, chatID)
}
