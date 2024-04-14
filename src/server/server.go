package main

import (
	"io"
	"log"
	"net"
	"time"
	"strconv"
	handler "user/ordersystem/src/handler"
	pb "user/ordersystem/src/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOrderManagementServer
}

func (s *server) GetOrderBidirectional(stream pb.OrderManagement_GetOrderBidirectionalServer) error {
	for {
			req, err := stream.Recv()
			if err == io.EOF {
					return nil
			}
			if err != nil {
					return err
			}
			found, orders := handler.FindOrderByItemName(req.Items)
			if found {
					for _, order := range orders {
							if err := stream.Send(&pb.OrderResponse{ItemName: order, TimeStamp: strconv.Itoa(time.Now().Second())}); err != nil {
									return err
							}
					}
			} else {
					if err := stream.Send(&pb.OrderResponse{ItemName: "No orders found for the specified item", TimeStamp: strconv.Itoa(time.Now().Second())}); err != nil {
							return err
					}
			}
	}
}


func (s *server) GetOrderServerStreaming(req *pb.OrderRequest, stream pb.OrderManagement_GetOrderServerStreamingServer) error {
	found, orders := handler.FindOrderByItemName(req.Items)
	if found {
		for _, order := range orders {
			if err := stream.Send(&pb.OrderResponse{ItemName: order, TimeStamp: strconv.Itoa(time.Now().Second())}); err != nil {
				return err
			} 
		}
	} else {
		if err := stream.Send(&pb.OrderResponse{ItemName: "No orders found for the specified item", TimeStamp: strconv.Itoa(time.Now().Second())}); err != nil {
				return err
		}
}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterOrderManagementServer(s, &server{})
	log.Println("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
