package common

import (
	"fmt"
	"log"
	"net"
	"os"

	RegistryServiceClient "github.com/sibaazab/large-scale-workshop.git/services/registry-service/client"
	//RegistryService "github.com/sibaazab/large-scale-workshop.git/services/registry-service/service"
	"gopkg.in/yaml.v2"

	"github.com/pebbe/zmq4"
	"github.com/sibaazab/large-scale-workshop.git/Config"
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

func LoadRegistryAddresses() []string {
	configFile := os.Args[1]
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading registry yaml file: %v", err)
		os.Exit(2)
	}
	var config Config.RegistryServiceConfig
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling registry addresses data: %v", err)
		os.Exit(3)
	}
	if len(config.RegistryAddresses) <= 0 {
		log.Fatalf("registry addresses yaml file does not include any enteries")
		os.Exit(4)
	}
	
	return config.RegistryAddresses
}


func registerAddress(serviceName string, registryAddresses []string, listeningAddress string) (unregister func()) {
    registryClient := RegistryServiceClient.NewRegistryServiceClient(registryAddresses)
    err := registryClient.Register(serviceName, listeningAddress)
    if err != nil {
            Logger.Fatalf("Failed to register to registry service: %v", err)
        }
    return func() {
    	registryClient.Unregister(serviceName, listeningAddress) }
    }

func Start(serviceName string, bindgRPCToService func(s grpc.ServiceRegistrar)) (int, func()) {
	listeningAddress, grpcServer, startListening, assignedPort := startgRPC()
	bindgRPCToService(grpcServer)
	
	
	unregister:= registerAddress(serviceName, LoadRegistryAddresses(), listeningAddress)
	//RegistryService.services[serviceName][listeningAddress]= unregister
	//go RegistryService.monitorNode(serviceName, listeningAddress)
	log.Printf("TestService listening on port %d", assignedPort)
	log.Printf("TestService listening on Address %v", listeningAddress)
	startListening()
	log.Print("After startlistening")
	return  assignedPort, unregister
}




func bindMQToService(listenPort int, messageHandler func(method string, parameters []byte) (response proto.Message, err error)) (startMQ func(), listeningAddress string){

	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
 		Logger.Fatalf("Failed to create a new zmq socket: %v", err)
	} 
	if listenPort == 0 {
 		listeningAddress = "tcp://127.0.0.1:*"
	} else {
 		listeningAddress = fmt.Sprintf("tcp://127.0.0.1:%v", listenPort)
	} 
	err = socket.Bind(listeningAddress)
	if err != nil {
 		Logger.Fatalf("Failed to bind a zmq socket: %v", err)
	} 
	listeningAddress, err = socket.GetLastEndpoint()
	if err != nil {
		 Logger.Fatalf("Failed to get listetning address of zmq socket: %v", err)
	} 

	startMQ = func() {
		for {
			data, readerr := socket.RecvBytes(0)
			if err != nil {
				Logger.Printf("Failed to receive bytes from MQ socket: %v\n", readerr)
				continue
			}
			if len(data) == 0 {
				continue
			}
			Logger.Printf("data len: %v\n", len(data))

			go func() {


			}
		}


	}
}
