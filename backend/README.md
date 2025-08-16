Käytä komentoa:
grpcurl -plaintext -proto proto/hello.proto -d '{\"name\":\"Paula\"}' localhost:50051 hello.Greeter/SayHello
