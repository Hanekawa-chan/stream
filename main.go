package main

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"stream/proto"
)

var upgrader = websocket.Upgrader{}

type Stream struct {
}

//func (s Stream) EchoStr(ctx context.Context, message *proto.Message) (*proto.Message, error) {
//	panic("implement me")
//}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Log().Err(err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Log().Err(err)
			break
		}
	}
}

func (s Stream) EchoReq(ctx context.Context, message *proto.Message) (*proto.Message, error) {
	if message != nil {
		return message, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "message came empty")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &Stream{}
	twirpHandler := proto.NewStreamServiceServer(server, nil)

	http.HandleFunc("/echo", echo)

	log.Log().Msg("working")

	go func() {
		if err := http.ListenAndServe(":8081", twirpHandler); err != nil {
			log.Fatal().Err(err)
		}
	}()
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal().Err(err)
	}
}
