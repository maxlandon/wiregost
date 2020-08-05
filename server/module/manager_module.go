package module

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

func (m *manager) RunModule(context.Context, *pb.RunRequest) (*pb.Run, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunModule not implemented")
}
func (m *manager) GetInfo(context.Context, *pb.InfoRequest) (*pb.Info, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (m *manager) GetOptions(context.Context, *pb.OptionsRequest) (*pb.Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOptions not implemented")
}
func (m *manager) SetOption(context.Context, *pb.SetOptionRequest) (*pb.Option, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOption not implemented")
}
func (m *manager) EditModule(context.Context, *pb.EditRequest) (*pb.Edit, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditModule not implemented")
}
func (m *manager) UseModule(context.Context, *pb.UseRequest) (*pb.Use, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseModule not implemented")
}
