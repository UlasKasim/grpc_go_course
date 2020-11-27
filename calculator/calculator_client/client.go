package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm calculator client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()
	defer fmt.Printf("\nClosed client")

	c := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Printf("Created client: %f", c)

	// doAddUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doAddUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Client sending two numbers 3 and 10...")
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			FirstValue:  3,
			SecondValue: 10,
		},
	}
	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while adding two numbers: %v", err)
	}
	log.Printf("Result of adding two numbers is: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Client handling decompose of primeNumber")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		PrimeNumber: 1381381312,
	}
	resStream, err := c.DecomposeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling DecomposeNumber RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from DecomposeNumber: %v", msg.GetDecompose())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			NumberToCompute: 131231,
		},
		{
			NumberToCompute: 1231322,
		},
		{
			NumberToCompute: 123123,
		},
		{
			NumberToCompute: 123123123,
		},
		{
			NumberToCompute: 12313131,
		},
		{
			NumberToCompute: 12313,
		},
		{
			NumberToCompute: 1213151,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage")
	}

	// We iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response %v", res)
}
