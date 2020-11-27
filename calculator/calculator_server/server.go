package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Add(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Add function was invoked with %v\n", req)
	firstValue := req.GetCalculator().GetFirstValue()
	secondValue := req.GetCalculator().GetSecondValue()
	res := &calculatorpb.CalculatorResponse{
		Result: firstValue + secondValue,
	}
	return res, nil
}

func (*server) DecomposeNumber(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_DecomposeNumberServer) error {
	fmt.Printf("DecomposeNumber function was invoked with %v\n", req)
	primeNumber := req.GetPrimeNumber()
	k := int32(2)
	for primeNumber > 1 {
		if primeNumber%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Decompose: k,
			}
			time.Sleep(1000 * time.Millisecond)
			stream.Send(res)
			primeNumber = primeNumber / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(reqStream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request %v\n", reqStream)
	result := int64(0)
	counter := int64(0)
	for {
		req, err := reqStream.Recv()
		if err == io.EOF {
			// we have finished the client stream
			return reqStream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				AverageOfNumbers: result / counter,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		counter++
		result += req.GetNumberToCompute()
	}
}

func (*server) FindMaximum(reqStream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request %v\n", reqStream.Context())
	resultNumber := int64(0)

	for {
		req, err := reqStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		number := req.GetNumber()
		if number > resultNumber {
			resultNumber = number
		}

		streamErr := reqStream.Send(&calculatorpb.FindMaximumResponse{
			Average: resultNumber,
		})
		if streamErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return streamErr
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
