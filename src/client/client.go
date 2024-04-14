package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	cls "user/ordersystem/pkg"
	pb "user/ordersystem/src/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetInputBidirectional() []string {
	var orders []string
	fmt.Println("Enter orders for Bidirectional Streaming:")
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

func ConnectToServer() (pb.OrderManagementClient, *grpc.ClientConn) {
	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	// defer conn.Close()
	client := pb.NewOrderManagementClient(conn)
	return client, conn
}

func RunClient(client pb.OrderManagementClient) {
	exit := false
	for {
		if exit {
			break
		}

		cls.CallClear()

		fmt.Println("Select Communication Pattern:\n1. Server Streaming\n2. Bidirectional Streaming\n3. Exit")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1":
			ServerStreaming(client)
		case "2":
			BidirectionalStreaming(client)
		case "3":
			exit = true
		default:
			fmt.Println("Invalid Input!")
		}
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
	}
}

func main() {

	client, connection := ConnectToServer()

	RunClient(client)

	connection.Close()
}
