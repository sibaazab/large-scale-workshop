package commonRegistry

import (
	//"context"
	//"log"
	"net"

	"fmt"
	"google.golang.org/grpc"
	. "github.com/sibaazab/large-scale-workshop.git/utils"
)

// startgRPC initializes and starts the gRPC server.
func startgRPC(listenPort int) (string, *grpc.Server, func(), error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		Logger.Printf("Failed to listen on port %v: %v", listenPort, err)
		return "", nil, nil, err
	}
	
	listeningAddress := lis.Addr().String()
	grpcServer := grpc.NewServer()

	startListening := func() {
		Logger.Printf("Starting gRPC server on %s", listeningAddress)
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Fatalf("Failed to serve gRPC: %v", err)
		}
	}

	return listeningAddress, grpcServer, startListening, nil
}

// Start initializes the server using configuration data and attempts to start it.
func Start(baseListenPort int, bindgRPCToService func(grpc.ServiceRegistrar)) (startListening func(),listenPort int, err error) {
	for i := 0; ; i++ {
		listenPort = baseListenPort + i
		listeningAddress:=""
		var grpcServer *grpc.Server
		listeningAddress, grpcServer, startListening, err = startgRPC(listenPort)
		if err != nil {
			continue
		}
		bindgRPCToService(grpcServer)
		Logger.Printf("RegistryService started at port %v, listening on %s", listenPort, listeningAddress)
		if listeningAddress!=""{
			break
		}
	}

	
	
	return startListening, listenPort, nil
}
