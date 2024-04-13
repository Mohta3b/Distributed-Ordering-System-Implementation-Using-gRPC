package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	pb "user/ordersystem/src/proto"

	"google.golang.org/grpc"
)

func GetInputBidirectional() []string {
	var orders []string
	fmt.Println("Enter orders for Bidirectional streaming:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}
		orders = append(orders, text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	return orders
}

func GetInputServerStreaming() string {
	fmt.Println("Enter order for Server streaming:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return scanner.Text()
}

func BidirectionalStreaming(client pb.OrderManagementClient) {
	orderRequests := GetInputBidirectional()
	getOrderClient, err := client.GetOrderBidirectional(context.Background())
	if err != nil {
		log.Fatalf("Error calling GetOrderBidirectional: %v", err)
	}
	for _, orderRequest := range orderRequests {
		request := &pb.OrderRequest{Items: orderRequest}
		if err := getOrderClient.Send(request); err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
	}

	for {
		orderResponse, err := getOrderClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}
		log.Printf("Order: %s", orderResponse)
	}
}

func ServerStreaming(client pb.OrderManagementClient) {
	orderRequest := &pb.OrderRequest{Items: GetInputServerStreaming()}
	getOrderClient, err := client.GetOrderServerStreaming(context.Background(), orderRequest)
	if err != nil {
		log.Fatalf("Error calling GetOrderServerStreaming: %v", err)
	}
	for {
		orderResponse, err := getOrderClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}
		log.Printf("Order: %s", orderResponse)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewOrderManagementClient(conn)
	fmt.Println("Enter inputs with space, press 'Enter' twice to finish")
	ServerStreaming(client)
	BidirectionalStreaming(client)
}
