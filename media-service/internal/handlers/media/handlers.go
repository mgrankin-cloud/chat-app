package media

import (
	"context"
	"errors"
	errors2 "github.com/mgrankin-cloud/messenger/pkg/storage"
	ssov6 "github.com/mgrankin-cloud/messenger/contract/gen/go/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov6.UnimplementedMediaServer
	media Media
}

type Media interface {
	UploadMedia(
		ctx context.Context,
		data []byte,
		fileName string,
		mimeType string,
	) (fileID int64, err error)
	DownloadMedia(
		ctx context.Context,
		fileID int64,
	) (data []byte, fileName string, mimeType string, err error)
}

func RegisterMediaService(gRPCServer *grpc.Server, media Media) {
	ssov6.RegisterMediaServer(gRPCServer, &serverAPI{media: media})
}

func (api *serverAPI) UploadMedia(
	ctx context.Context,
	req *ssov6.UploadMediaRequest,
) (*ssov6.UploadMediaResponse, error) {

	if req.GetData() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "file data is required")
	}

	fileID, err := api.media.UploadMedia(ctx, req.GetData(), req.GetFileName(), req.GetMimeType())
	if err != nil {
		if errors.Is(err, errors2.ErrFileExists) {
			return nil, status.Error(codes.AlreadyExists, "file already exists")
		}

		return nil, status.Error(codes.Internal, "failed to upload file")
	}
	
	return &ssov6.UploadMediaResponse {
		FileId: fileID,
	}, nil
}

func (api *serverAPI) DownloadMedia(
	ctx context.Context,
	req *ssov6.DownloadMediaRequest,
) (*ssov6.DownloadMediaResponse, error) {

	if req.GetFileId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "file id is required")
	}

	data, fileName , mimeType, err := api.media.DownloadMedia(ctx, req.GetFileId())
	if err != nil {
		if errors.Is(err, errors2.ErrFileNotFound) {
			return nil, status.Error(codes.AlreadyExists, "file not found")
		}
		return nil, status.Error(codes.Internal, "failed to download file")
	}

	return &ssov6.DownloadMediaResponse{
		Data: data,
		FileName: fileName,
		MimeType: mimeType,
	}, nil
}
