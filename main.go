package main

import (
	"context"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"stream/proto"
	"time"
)

type Stream struct {
	proto.UnimplementedStreamServiceServer
}

func (s *Stream) EchoReq(_ context.Context, empty *proto.Empty) (*proto.Result, error) {
	log.Log().Msg("get " + empty.Msg)
	return &proto.Result{Status: true, Msg: "ok"}, nil
}

func (s *Stream) EchoStr(empty *proto.Empty, client proto.StreamService_EchoStrServer) error {
	for {
		time.Sleep(10 * time.Second)
		log.Log().Msg("get " + empty.Msg)
		err := client.Send(&proto.Result{Status: true, Msg: "sending something"})
		if err != nil {
			return err
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lis, err := net.Listen("tcp", ":8842")
	if err != nil {
		log.Fatal().Err(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStreamServiceServer(grpcServer, &Stream{})
	grpc_prometheus.Register(grpcServer)

	var group errgroup.Group

	group.Go(func() error {
		return grpcServer.Serve(lis)
	})

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	runtime.SetHTTPBodyMarshaler(gwmux)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}

	group.Go(func() error {
		return proto.RegisterStreamServiceHandlerFromEndpoint(ctx, gwmux, ":8842", opts)
	})
	mux := http.NewServeMux()
	mux.Handle("/", wsproxy.WebsocketProxy(gwmux))
	group.Go(func() error {
		return http.ListenAndServe(":8843", mux)
	})
	group.Go(func() error {
		return http.ListenAndServe(":8844", promhttp.Handler())
	})
	log.Log().Msg("working")
	err = group.Wait()
	if err != nil {
		log.Fatal().Err(err)
	}
}
