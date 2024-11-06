package main

import "grpc-file-service/cmd/server"

// main - The main entry point for the gRPC file service.
//
// This main function is used to start the gRPC server and
// listen for incoming requests.
func main() {
	// Run the gRPC server
	server.Run()
}
