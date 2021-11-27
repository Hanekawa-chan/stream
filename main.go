package main

import (
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog/log"
	"github.com/tmc/grpc-websocket-proxy/examples/cmd/wsechoserver/echoserver"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"stream/proto"
	"time"
)

type Stream struct {
	proto.UnimplementedStreamServiceServer
}

func (s *Stream) EchoReq(ctx context.Context, empty *proto.Empty) (*proto.Result, error) {
	return &proto.Result{Status: true, Msg: "ok"}, nil
}

func(s *Stream) EchoStr(_ *proto.Empty, client proto.StreamService_EchoStrServer) error {
	for {
		time.Sleep(10*time.Second)
		err := client.Send(&proto.Result{Status: true, Msg: "sending something"})
		if err != nil {
			return err
		}
	}
}

func main() {
	server := grpc.NewServer()
	proto.RegisterStreamServiceServer(server, &Stream{})
	grpcPrometheus.Register(server)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	go func() {
		if err := echoserver.RegisterEchoServiceHandlerFromEndpoint(context.Background(), mux, ":9998", opts); err != nil {
			log.Fatal().Err(err)
		}
		err := http.ListenAndServe(":10000", wsproxy.WebsocketProxy(mux))
		if err != nil {
			log.Fatal().Err(err)
		}
	}()

	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	log.Log().Msg("public server started")
	if err := server.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("listen server")
	}
}