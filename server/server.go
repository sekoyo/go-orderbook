package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/lightfin-io/orderbook/orderbook"
	pb "github.com/lightfin-io/orderbook/proto"
	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	baseCcy  = flag.String("base_ccy", "BTC", "The base currency e.g. BTC")
	quoteCcy = flag.String("quote_ccy", "USD", "The qoute currency e.g. USD")
)

type server struct {
	ob *orderbook.Orderbook
	pb.UnimplementedOrderbookServer
}

func (s *server) AddOrder(ctx context.Context, in *pb.AddOrderRequest) (*pb.AddOrderReply, error) {
	log.Printf("Received AddOrder")
	var side orderbook.OrderSide
	if in.IsBid {
		side = orderbook.Bid
	} else {
		side = orderbook.Ask
	}
	err := s.ob.AddOrder(in.OrderId, side, in.Price, in.Qty)
	return &pb.AddOrderReply{Error: err.Error()}, nil
}

func (s *server) CancelOrder(ctx context.Context, in *pb.CancelOrderRequest) (*pb.CancelOrderReply, error) {
	log.Printf("Received CancelOrder %d", in.OrderId)
	err := s.ob.CancelOrder(in.OrderId)
	return &pb.CancelOrderReply{Error: err.Error()}, nil
}

func (s *server) AmendOrder(ctx context.Context, in *pb.AmendOrderRequest) (*pb.AmendOrderReply, error) {
	log.Printf("Received AmendOrder %d", in.OrderId)
	err := s.ob.AmendOrder(in.OrderId, in.Qty)
	return &pb.AmendOrderReply{Error: err.Error()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ob := orderbook.NewOrderbook(*baseCcy, *quoteCcy)
	s := grpc.NewServer()
	pb.RegisterOrderbookServer(s, &server{ob: ob})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
