package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/imaskm/greet/constants"
	"github.com/imaskm/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fname := req.GetGreeting().GetFirstName()

	result := "Hello " + fname

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil

}

func (*server) GreetMany(req *greetpb.GreetManyRequest, stream greetpb.GreetService_GreetManyServer) error {

	name := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "hello " + name + " times " + strconv.Itoa(i+1)
		stream.Send(&greetpb.GreetManyResponse{
			Result: result,
		})
		time.Sleep(time.Millisecond * 1000)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Streaming ended!!")
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatal(err)
		}
		name := msg.GetGreeting().GetFirstName()
		log.Println("Received name: ", name)
		result += "Hello " + name + "!! "
	}
}

func errCheck(ctx context.Context) {
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Client canceled the request")
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {

	defer errCheck(ctx)
	time.Sleep(3 * time.Second)

	fname := req.GetGreeting().GetFirstName()

	result := "Hello Deadline " + fname

	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil

}

func main() {
	fmt.Println("Hello World")

	listener, err := net.Listen("tcp", constants.Server+":"+constants.Port)
	if err != nil {
		log.Fatal(err)
	}

	tls := false
	opts := []grpc.ServerOption{}
	if tls {
		c_path := "./ssl/"

		cred, err := credentials.NewServerTLSFromFile(c_path+"server.crt", c_path+"server.pem")

		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, grpc.Creds(cred))
	}

	s := grpc.NewServer(opts...)

	reflection.Register(s)

	go func() {
		greetpb.RegisterGreetServiceServer(s, &server{})
		log.Println("Started the server on:  ", constants.Server+":"+constants.Port)
		if err = s.Serve(listener); err != nil {
			log.Fatal("failed", err)
		}
	}()

	// conn, err := grpc.DialContext(
	// 	context.Background(),
	// 	"0.0.0.0:5000",
	// 	grpc.WithBlock(),
	// 	grpc.WithInsecure(),
	// )
	// if err != nil {
	// 	log.Fatalln("Failed to dial server:", err)
	// }

	// gwmux := runtime.NewServeMux()
	// // Register Greeter
	// err = greetpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	// if err != nil {
	// 	log.Fatalln("Failed to register gateway:", err)
	// }

	// gwServer := &http.Server{
	// 	Addr:    ":8090",
	// 	Handler: gwmux,
	// }

	// log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	// log.Fatalln(gwServer.ListenAndServe())
}
