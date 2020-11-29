package main

import (
	"context"
	"fmt"
	"grpc-go-course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create Test Blog
	blog := &blogpb.Blog{
		Title:    "Ipek Naber",
		Content:  "İyidir senden naber",
		AuthorId: "RandomId",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("An error occured during adding blog with %v", err)
	}
	fmt.Printf("Blog has been created: %v", res)
}
