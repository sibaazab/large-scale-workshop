package RegistryService

import (
	"context"
	"log"
	//"net"

	//"fmt"
	//"sync"
	//"time"
	commonRegistry "github.com/sibaazab/large-scale-workshop.git/services/commonRegistry"
	. "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
	RegistryServiceServant "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"

	"github.com/sibaazab/large-scale-workshop.git/Config"
	. "github.com/sibaazab/large-scale-workshop.git/utils"
)


type registeryServiceImplementation struct {
	UnimplementedRegistryServiceServer
}


// Start initializes the server using configuration data and attempts to start it.
func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) { 
        RegisterRegistryServiceServer(s, &registeryServiceImplementation{})
    }
	var config Config.RegistryConfigBase
	if err := yaml.Unmarshal(configData, &config); err != nil {
		Logger.Fatalf("Error unmarshaling configuration data: %v", err)
		return err
	}

	baseListenPort := config.ListenPort
	log.Printf("----------------baselistenPort=%v", baseListenPort)
    startListening, listenPort, err := commonRegistry.Start(baseListenPort, bindgRPCToService)
	if err != nil {
		Logger.Printf("Failed to start the Registry service %v: ", err)
		return  nil
	}
	log.Printf("-------------------base=%v listen=%v", baseListenPort, listenPort)
	if err := RegistryServiceServant.CreateChord(baseListenPort, listenPort); err != nil {
		Logger.Printf("Error creating chord: %v", err)
		return err
	}
	startListening()
	//defer unregister()
    return nil
}

func (obj *registeryServiceImplementation) Register(ctxt context.Context, req *ServiceRequest) (e *emptypb.Empty, err error) {
	// //servant := NewRegistryService()
	log.Printf("sdfghjkllkjhgfd0")
	err1 := RegistryServiceServant.Register(req.GetName(), req.GetAddress())
	if err1 != nil {
		log.Printf("Failed to register service: %v", err)
		return nil, err1
	}
	
	return &emptypb.Empty{}, nil
}


func (obj *registeryServiceImplementation) Unregister(ctxt context.Context, req *ServiceRequest) (e *emptypb.Empty, err error) {
	
	err1 := RegistryServiceServant.Unregister(req.GetName(), req.GetAddress())
	if err1 != nil {
		log.Printf("Failed to unregister service: %v", err1)
		return nil, err1
	}
	
	return &emptypb.Empty{}, nil
}

func (obj *registeryServiceImplementation) Discover(ctxt context.Context, in *wrapperspb.StringValue) (SNode *ServiceNodes, err error) {
	//servant := NewRegistryService()
	addresses := RegistryServiceServant.Discover(in.GetValue())
	if len(addresses) == 0 {
		return &ServiceNodes{Nodes: []string{}}, nil
	}
	return &ServiceNodes{Nodes: []string{}}, nil
}
