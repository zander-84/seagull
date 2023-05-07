package grpc_router

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/transport"
	grpc2 "github.com/zander-84/seagull/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Router struct {
	rmc           endpoint.Rmc
	serverHandler []grpc.ServiceDesc
}

func NewRouter(rmc endpoint.Rmc) *Router {
	router := new(Router)
	router.rmc = rmc
	router.serverHandler = make([]grpc.ServiceDesc, 0)
	return router
}
func (r *Router) ServerHandler() []grpc.ServiceDesc {
	return r.serverHandler
}

func grpcCtx(ctx context.Context, fullPath endpoint.Path) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	endpointCtxVal := transport.NewTransporter(transport.Grpc, transport.NewHeader(md), fullPath.FullPath(), fullPath.Method())
	ctx = transport.WithContext(ctx, endpointCtxVal)
	return grpc2.NewGrpcContext(ctx, endpointCtxVal)
}

func (r *Router) Endpoint(kind transport.Kind, fullPath endpoint.Path, h endpoint.HandlerFunc) {
	data := grpc.ServiceDesc{
		ServiceName: fullPath.ServerName(),
		HandlerType: nil,
		Methods: []grpc.MethodDesc{
			{
				MethodName: fullPath.Method().String(),
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

					//  grpc默认先dec，gull把dec部分后置，在endpoint前实现
					in := dec
					//codec, err := r.rmc.GetCodec(protocol, method, path)
					//if err != nil {
					//	return nil, err
					//}
					//in := codec.NewReq()
					//if err := dec(in); err != nil {
					//	return nil, err
					//}

					if interceptor == nil {
						return h(grpcCtx(ctx, fullPath), in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: fullPath.FullPath(),
					}
					handler2 := func(ctx context.Context, req interface{}) (interface{}, error) {
						return h(grpcCtx(ctx, fullPath), in)
					}
					return interceptor(ctx, in, info, handler2)
				},
			},
		},
		Streams:  nil,
		Metadata: nil,
	}
	r.serverHandler = append(r.serverHandler, data)
}
