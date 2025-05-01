package main

import (
	"fmt"

	"context"
	example "github.com/micro/micro/examples/server/proto/example"
	"go-micro.dev/v5"
)

// publishes a message
func pub(i int, p micro.Publisher) {
	msg := &example.Message{
		Say: fmt.Sprintf("This is an async message %d", i),
	}

	if err := p.Publish(context.TODO(), msg); err != nil {
		fmt.Println("pub err: ", err)
		return
	}

	fmt.Printf("Published %d: %v\n", i, msg)
}

func main() {
	service := micro.NewService()
	service.Init()

	p := micro.NewEvent("example", service.Client())

	fmt.Println("\n--- Publisher example ---")

	for i := 0; i < 10; i++ {
		pub(i, p)
	}
}
