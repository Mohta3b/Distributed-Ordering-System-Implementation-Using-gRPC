package main

// import "google.golang.org/grpc"
import (
	"fmt"
	"user/ordersystem/src/handler"
)

func PrintThat() string {
	return "hello from server"
}

func main() {
	fmt.Println(PrintThat())
	fmt.Println(handler.Alaki())
}
