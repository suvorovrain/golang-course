





package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)




const _ = grpc.SupportPackageIsVersion9

const (
	CollectorService_GetRepoInfo_FullMethodName = "/collector.CollectorService/GetRepoInfo"
)




type CollectorServiceClient interface {
	GetRepoInfo(ctx context.Context, in *RepoRequest, opts ...grpc.CallOption) (*RepoResponse, error)
}

type collectorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCollectorServiceClient(cc grpc.ClientConnInterface) CollectorServiceClient {
	return &collectorServiceClient{cc}
}

func (c *collectorServiceClient) GetRepoInfo(ctx context.Context, in *RepoRequest, opts ...grpc.CallOption) (*RepoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoResponse)
	err := c.cc.Invoke(ctx, CollectorService_GetRepoInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}




type CollectorServiceServer interface {
	GetRepoInfo(context.Context, *RepoRequest) (*RepoResponse, error)
	mustEmbedUnimplementedCollectorServiceServer()
}






type UnimplementedCollectorServiceServer struct{}

func (UnimplementedCollectorServiceServer) GetRepoInfo(context.Context, *RepoRequest) (*RepoResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetRepoInfo not implemented")
}
func (UnimplementedCollectorServiceServer) mustEmbedUnimplementedCollectorServiceServer() {}
func (UnimplementedCollectorServiceServer) testEmbeddedByValue()                          {}




type UnsafeCollectorServiceServer interface {
	mustEmbedUnimplementedCollectorServiceServer()
}

func RegisterCollectorServiceServer(s grpc.ServiceRegistrar, srv CollectorServiceServer) {
	
	
	
	
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CollectorService_ServiceDesc, srv)
}

func _CollectorService_GetRepoInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RepoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectorServiceServer).GetRepoInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectorService_GetRepoInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectorServiceServer).GetRepoInfo(ctx, req.(*RepoRequest))
	}
	return interceptor(ctx, in, info, handler)
}




var CollectorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "collector.CollectorService",
	HandlerType: (*CollectorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRepoInfo",
			Handler:    _CollectorService_GetRepoInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/collector.proto",
}
