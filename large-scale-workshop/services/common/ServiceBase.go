package services

import (
	//"fmt"
	"log"
	"net"

	RegistryServiceClient "github.com/sibaazab/large-scale-workshop.git/services/registry-service/client"

	. "github.com/sibaazab/large-scale-workshop.git/utils"
	"google.golang.org/grpc"
)

func startgRPC() (listeningAddress string, grpcServer *grpc.Server,
	startListening func(), assignedPort int) {
	lis, err := net.Listen("tcp", ":0")
	assignedPort = lis.Addr().(*net.TCPAddr).Port
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		Logger.Fatalf("failed to listen(ServiceBase): %v", err)
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

func Start(serviceName string, bindgRPCToService func(s grpc.ServiceRegistrar)) int {
	listeningAddress, grpcServer, startListening, assignedPort := startgRPC()
	bindgRPCToService(grpcServer)
	log.Printf("TestService listening on port %d", assignedPort)
	log.Printf("TestService listening on Address %v", listeningAddress)
	startListening()
	log.Print("After startlistening")
	return assignedPort
}

func registerAddress(serviceName string, registryAddresses []string, listeningAddress string) (unregister func()) {
	registryClient := RegistryServiceClient.NewRegistryServiceClient(registryAddresses)
	err := registryClient.Register(serviceName, listeningAddress)
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	return func() {
		registryClient.Unregister(serviceName, listeningAddress)
	}
}
