package api

import (
	"context"

	pb "auth.service/api/proto"
	"auth.service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccessServiceHandler struct {
	pb.UnimplementedAccessServiceServer
	accessService service.AcccessService
}

func NewAccessServiceHandler(accessService service.AcccessService) *AccessServiceHandler {
	return &AccessServiceHandler{
		accessService: accessService,
	}
}

func (h *AccessServiceHandler) Check(
	ctx context.Context,
	req *pb.CheckAccessRequest,
) (*pb.CheckAccessResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"access token is empty",
		)
	}

	isValid, userID, err := h.accessService.Check(ctx, req.AccessToken)
	if err != nil {
		return &pb.CheckAccessResponse{
			IsValid: false,
			UserId:  "",
		}, nil
	}

	return &pb.CheckAccessResponse{
		IsValid: isValid,
		UserId:  userID,
	}, nil
}
