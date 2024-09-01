package auth

import (
	ssov6 "Messenger-android/messenger/media-service/contract/gen/go/media"
	"context"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov6.UnimplementedMediaServer
	media Media
}

type Media interface {
	UploadMedia(
		ctx context.Context,
		req *ssov6.UploadMediaRequest,
	) (ssov6.UploadMediaResponse, error)
	DownloadMedia(
		ctx context.Context,
		req *ssov6.DownloadMediaRequest,
	) (ssov6.DownloadMediaResponse, error)
}

func RegisterMediaServer(gRPCServer *grpc.Server, media Media) {
	ssov6.RegisterMediaServer(gRPCServer, &serverAPI{media: media})
}

func (api *serverAPI) UploadMedia(
	ctx context.Context,
	req *ssov6.UploadMediaRequest,
) (ssov6.UploadMediaResponse, error) {
	return api.media.UploadMedia(ctx, req)
}

func (api *serverAPI) DownloadMedia(
	ctx context.Context,
	req *ssov6.DownloadMediaRequest,
) (ssov6.DownloadMediaResponse, error) {
	return api.media.DownloadMedia(ctx, req)
}
