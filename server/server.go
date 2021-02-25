package server

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*Server) Decomposition(req *calculatorpb.DecompositionRequest, stream calculatorpb.CalculatorService_DecompositionServer) error {
	fmt.Printf("Decomposition function was invoked with %v \n", req)
	number := req.GetNum()
	var divide int64 = 2

	for number > 1 {
		if number%divide == 0 {
			res := &calculatorpb.DecompositionResponse{Decompose: divide}
			err := stream.Send(res)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
			number = number / divide
			time.Sleep(time.Second)
		} else {
			divide++
		}
	}

	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called..")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Avg: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("err while Recv Average %v", err)
			return err
		}
		log.Printf("receive req %v", req)
		total += req.GetNum()
		count++
	}
}
func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
