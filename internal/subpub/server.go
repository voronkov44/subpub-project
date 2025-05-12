package subpub

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"subpub-project/proto"
)

type server struct {
	pubSub                          PubSub
	proto.UnimplementedPubSubServer // встраиваем UnimplementedPubSubServer
}

func (s *server) MustEmbedUnimplementedPubSubServer() {}

func (s *server) Subscribe(req *proto.SubscribeRequest, stream proto.PubSub_SubscribeServer) error {
	sub, err := s.pubSub.Subscribe(req.Key, func(msg interface{}) {
		// Отправка сообщений подписчику
		event := &proto.Event{
			Data: fmt.Sprintf("%v", msg),
		}
		if err := stream.Send(event); err != nil {
			log.Println("Error sending event:", err)
		}
	})
	if err != nil {
		return err
	}
	// Ожидаем завершения подписки
	<-stream.Context().Done()
	sub.Unsubscribe()
	return nil
}

func (s *server) Publish(ctx context.Context, req *proto.PublishRequest) (*emptypb.Empty, error) { // Используем emptypb.Empty
	if err := s.pubSub.Publish(req.Key, req.Data); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil // возвращаем пустую структуру
}

func StartGRPCServer(address string, pubSub PubSub) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterPubSubServer(s, &server{pubSub: pubSub})

	reflection.Register(s)

	log.Println("Starting gRPC server on", address)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
