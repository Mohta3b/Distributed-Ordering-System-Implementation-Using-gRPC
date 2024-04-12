package main

// import "google.golang.org/grpc"

import (
	"user/ordersystem/src/handler"
)

func PrintThat() string {
	return "hello from server"
}

func main()	{
	PrintThat()
	handler.Alaki()
}