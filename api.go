package main

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc/status"

	"github.com/dxvgef/tsing-gateway/api"

	gw "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 启动API GRPC服务
func startApiGrpcServer() {
	// go apiEngine.start(proxyEngine)
	// 创建一个grpc服务的实例
	svr := grpc.NewServer()

	// 监听一个地址
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	// 注册user服务
	api.RegisterAPIServer(svr, &api.SourceHandler{})

	// 启用反射服务，允许客户端查询本实例提供的服务和方法
	reflection.Register(svr)

	// 启动grpc服务
	log.Info().Msg("启动GRPC Server")
	if err = svr.Serve(lis); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
}

// 启动API GRPC Gateway服务
func startApiGrpcGatewayServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建一个grpc gateway服务实例
	mux := gw.NewServeMux(
		// 自定义grpc的错误处理
		gw.WithProtoErrorHandler(errorHandler),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// 注册user服务的端点
	err := api.RegisterAPIHandlerFromEndpoint(ctx, mux, ":3000", opts)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	// 启动grpc gateway服务
	log.Info().Msg("启动GRPC Gateway Server")
	if err = http.ListenAndServe(":13000", mux); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
}

type ResponseError struct {
	Error string `json:"error,omitempty"`
	// Code  int32  `json:"code,omitempty"`
}

// 错误响应处理器
func errorHandler(_ context.Context, _ *gw.ServeMux, marshaller gw.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, _ := status.FromError(err)

	// // 如果不是400错误，则不返回任何消息，可以减少json序列化的过程以及网络传输
	// // 客户端可以根据状态码来决定输出什么消息
	// if s.Code() != 400 {
	// 	w.WriteHeader(gw.HTTPStatusFromCode(s.Code()))
	// 	return
	// }

	var respErr ResponseError
	// respErr.Code = s.Proto().GetCode()
	respErr.Error = s.Proto().GetMessage()

	jsonData, err := marshaller.Marshal(&respErr)
	if err != nil {
		w.WriteHeader(int(s.Code()))
		return
	}

	w.WriteHeader(gw.HTTPStatusFromCode(s.Code()))
	_, err = w.Write(jsonData)
	if err != nil {
		w.WriteHeader(int(s.Code()))
		return
	}
}
