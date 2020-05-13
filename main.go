package main

import (
	"context"
	"flag"
	"net"
	"net/http"

	gw "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/dxvgef/tsing-gateway/api"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/source"
)

func main() {
	// 设置默认logger
	setDefaultLogger()

	// 加载配置文件
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	err := global.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	// 获得一个代理引擎实例
	proxyEngine := proxy.New()

	// 启动api服务
	if global.Config.API.On {
		// 启动grpc服务
		go func() {
			// go apiEngine.Start(proxyEngine)
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
		}()
		// 启用grpc gateway服务
		go func() {
			// ------------------------- 启用grpc gateway server ---------------------------
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
		}()
	}

	// 启动网关引擎
	proxyEngine.Start()
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

// 初始化数据，目前仅开发调试用途
func initData(e *proxy.Engine, dataSource source.Source) (err error) {
	var (
		upstream   proxy.Upstream
		routeGroup proxy.RouteGroup
	)
	upstream.ID = "testUpstream"
	upstream.Middleware = append(upstream.Middleware, proxy.Configurator{
		Name:   "favicon",
		Config: `{"status": 204}`,
	})
	upstream.Discover.Name = "coredns_etcd"
	upstream.Discover.Config = `{"host":"test.uam.local"}`
	// 设置上游
	err = e.NewUpstream(upstream, false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	// 设置上游
	upstream = proxy.Upstream{}
	upstream.ID = "test2Upstream"
	upstream.Discover.Name = "coredns_etcd"
	upstream.Discover.Config = `{"host":"test2.uam.local"}`
	err = e.NewUpstream(upstream, false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	// 设置路由组
	routeGroup, err = e.SetRouteGroup("testGroup", false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}
	// 设置路由
	err = routeGroup.SetRoute("/test", "get", "testUpstream", false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}
	// 设置主机
	err = e.SetHost("127.0.0.1", routeGroup.ID, false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	// 序列化成json
	log.Debug().Interface("配置", e).Send()

	// 将所有数据保存到数据源
	if err = dataSource.SaveAll(); err != nil {
		log.Err(err).Caller().Send()
		return
	}

	return nil
}
