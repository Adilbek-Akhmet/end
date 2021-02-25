package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func doManyTimesFromServer(c calculatorpb.CalculatorServiceClient) {
	ctx := context.Background()
	req := &calculatorpb.DecompositionRequest{Num: 120}

	stream, err := c.Decomposition(ctx, req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break LOOP
		}
		if err != nil {
			log.Fatalf("error %v", err)
		}
		log.Printf("response:%v \n", res.GetDecompose())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Average server streaming RPC...")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	numbers := []int32{3, 8, 12}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.AverageRequest{
			Num: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	fmt.Printf("The Average is: %v\n", res.GetAvg())
}


	ctx := context.Background()
	stream, err := c.Average(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Average Response: %v\n", res)
}

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	doLongCalculateAverage(c)
}
