package services

import (
	"fmt"
	"log"
	"net"

	. "github.com/sibaazab/large-scale-workshop.git/utils"
	"google.golang.org/grpc"
) 
  
func startgRPC(listenPort int) (listeningAddress string, grpcServer *grpc.Server, 
startListening func()) { 
    lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort)) 
    if err != nil { 
        Logger.Fatalf("failed to listen: %v", err) 
    } 
    listeningAddress = lis.Addr().String() 
    grpcServer = grpc.NewServer() 
    startListening = func() { 
        if err := grpcServer.Serve(lis); err != nil { 
            Logger.Fatalf("failed to serve: %v", err) 
        } 
    } 
    return 
} 
  
func Start(serviceName string, grpcListenPort int, bindgRPCToService func(s grpc.ServiceRegistrar)) { 
    _ , grpcServer, startListening := startgRPC(grpcListenPort) // listeningAddress
    bindgRPCToService(grpcServer) 
    startListening()
    log.Print("After startlistening") 
} 
