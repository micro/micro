package handler

import (
	"context"

	pb "github.com/micro/go-micro/auth/service/proto"
)

// New returns an instance of Handler
func New() *Handler {
	return new(Handler)
}

type Handler struct{}

func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	return nil
}
func (h *Handler) Revoke(ctx context.Context, req *pb.RevokeRequest, rsp *pb.RevokeResponse) error {
	return nil
}
func (h *Handler) AddRole(ctx context.Context, req *pb.AddRoleRequest, rsp *pb.AddRoleResponse) error {
	return nil
}
func (h *Handler) RemoveRole(ctx context.Context, req *pb.RemoveRoleRequest, rsp *pb.RemoveRoleResponse) error {
	return nil
}
