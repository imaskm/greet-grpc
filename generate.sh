#protoc greetpb/greet.proto --go_out=.  --go-grpc_out=.
protoc greetpb/greet.proto --go_out=.  --go-grpc_out=require_unimplemented_servers=false:.