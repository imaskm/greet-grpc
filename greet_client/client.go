package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	
	"github.com/imaskm/greet/constants"
	"github.com/imaskm/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client started")
	tls := false
	opts := grpc.WithInsecure()

	if tls {
		creds, err := credentials.NewClientTLSFromFile("./ssl/ca.crt", "")

		if err != nil {
			log.Fatal(err)
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	grpc_client, err := grpc.Dial(constants.Server+":"+constants.Port, opts)

	if err != nil {
		log.Fatal("client fataa: ", err)
	}

	defer grpc_client.Close()
	c := greetpb.NewGreetServiceClient(grpc_client)
	doUnary(&c)

	// doServerStreaming(c)
	// doClientStreaming(c)
	// doUnaryWithDeadline(c, time.Second*5)
	// doUnaryWithDeadline(c, time.Second*1)

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	clientLong, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	namesArr := []string{"Ashwani Kaushik", "Sachin Tendulakr", "Nitin Goyal"}

	for _, name := range namesArr {
		names := strings.Split(name, " ")
		clientLong.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:  names[0],
				SecondName: names[1],
			},
		})
	}

	serverResp, err := clientLong.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response from server: ", serverResp.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {

	resStream, err := c.GreetMany(context.Background(), &greetpb.GreetManyRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Ashwani",
			SecondName: "Kaushik",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			log.Println("Streaming ended..")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Response from stream: ", msg.Result)
	}

}

func doUnary(c *greetpb.GreetServiceClient) {
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Ashwani",
			SecondName: "Kaushik",
		},
	}
	resp, err := (*c).Greet(context.Background(), request)

	if err != nil {
		log.Fatal("failed to call unary API:  ", err.Error())
	}

	log.Printf("Response from greet : %v", resp.Result)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName:  "Ashwani",
			SecondName: "Kaushik",
		},
	})

	if err != nil {
		s, ok := status.FromError(err)

		if ok && s.Code() == codes.DeadlineExceeded {
			fmt.Println("Deadline exceeded")
		} else {
			log.Fatal(err)
		}

		return
	}

	fmt.Println(res.GetResult())

}
