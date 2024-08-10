package ResgistryService

import (
	"context"
	"net"

	"fmt"
	//"sync"
	//"time"
	. "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
	"gopkg.in/yaml.v2"
	RegistryServiceServant "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant"
	servant "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	. "github.com/sibaazab/large-scale-workshop.git/utils"
	"github.com/sibaazab/large-scale-workshop.git/Config"
)


type registeryServiceImplementation struct {
	UnimplementedRegistryServiceServer
}

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

// StartServer sets up and starts the server with gRPC.
func helperStart(baseGrpcListenPort int, grpcListenPort int, bindgRPCToService func(grpc.ServiceRegistrar)) error {
	listeningAddress, grpcServer, startListening, err := startgRPC(grpcListenPort)
	if err != nil {
		return err
	}

	bindgRPCToService(grpcServer)
	Logger.Printf("RegistryService started at port %v, listening on %s", grpcListenPort, listeningAddress)

	// Create chord and check alive logic
	if err := RegistryServiceServant.CreateChord(baseGrpcListenPort, grpcListenPort); err != nil {
		Logger.Printf("Error creating chord: %v", err)
		return err
	}

	// if grpcListenPort == baseGrpcListenPort {
	// 	go RegistryServiceServant.CheckIsAliveEvery10Seconds()
	// }

	startListening()
	return nil
}

// Start initializes the server using configuration data and attempts to start it.
func Start(configData []byte) error {
	var config Config.RegistryConfigBase
	if err := yaml.Unmarshal(configData, &config); err != nil {
		Logger.Fatalf("Error unmarshaling configuration data: %v", err)
		return err
	}

	baseListenPort := config.ListenPort
	var startError error

	for i := 0; ; i++ {
		listenPort := baseListenPort + i
		if err := helperStart(baseListenPort, listenPort, func(s grpc.ServiceRegistrar) {
			RegisterRegistryServiceServer(s, &registeryServiceImplementation{})
		}); 
		err == nil {
			break
		} else {
			startError = err
			Logger.Printf("Error starting server on port %v: %v", listenPort, err)
		}
	}

	return startError
}

func (obj *registeryServiceImplementation) Register(ctxt context.Context, sv *ServiceRequest) (e *emptypb.Empty, err error) {
	//servant := NewRegistryService()
	//log.Printf("sdfghjkllkjhgfd0")
	servant.Register(sv.GetName(), sv.GetAddress())
	return
}

func (obj *registeryServiceImplementation) unregister(ctxt context.Context, sv *ServiceRequest) (e *emptypb.Empty, err error) {
	//servant := NewRegistryService()
	//servant= make(map[string]map[string]int)
	servant.Unregister(sv.GetName(), sv.GetAddress())
	return
}

func (obj *registeryServiceImplementation) Discover(ctxt context.Context, in *wrapperspb.StringValue) (SNode *ServiceNodes, err error) {
	//servant := NewRegistryService()
	addresses := servant.Discover(in.GetValue())
	if len(addresses) == 0 {
		return &ServiceNodes{Nodes: []string{}}, nil
	}

	return &ServiceNodes{Nodes: addresses}, nil
}
