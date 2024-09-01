package auth

import (
	ssov7 "Messenger-android/messenger/notification-service/contract/gen/go/notification"
	"context"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov7.UnimplementedNotificationServer
	notification Notification
}

type Notification interface {
	SendNotification(
		ctx context.Context,
		in *ssov7.SendNotificationRequest,
	) (*ssov7.SendNotificationResponse, error)
	SubscribeNotification(
		ctx context.Context,
		in *ssov7.SubscribeNotificationsRequest,
	) (*ssov7.SendNotificationResponse, error)
}

func RegisterNtfServer(gRPCServer *grpc.Server, notification Notification) {
	ssov7.RegisterNtfServer(gRPCServer, &serverAPI{notification: notification})
}

func (s *serverAPI) SendNotification(ctx context.Context, req *ssov7.SendNotificationRequest) (*ssov7.SendNotificationResponse, error) {
	message := req.GetMessage()
	userID := req.GetUserId()

	if subscriber, ok := s.subscribers.Load(userID); ok {
		subscriber.(chan *ssov7.Notification) <- &ssov7.Notification{Message: message}
	}

	return &ssov7.SendNotificationResponse{Success: true, Message: "Notification sent"}, nil
}

func (s *serverAPI) SubscribeNotifications(req *ssov7.SubscribeNotificationsRequest, stream ssov7.NotificationService_SubscribeNotificationsServer) error {
	userID := req.GetUserId()
	ch := make(chan *ssov7.Notification)
	s.subscribers.Store(userID, ch)

	defer func() {
		s.subscribers.Delete(userID)
		close(ch)
	}()

	for {
		select {
		case notification := <-ch:
			if err := stream.Send(notification); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
