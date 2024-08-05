package ResgistryService

import (
	"context"
	"log"
	"net"

	//"fmt"
	//"sync"
	//"time"
	. "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
	//RegistryServiceServant "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant"
	servant "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	//. "github.com/sibaazab/large-scale-workshop.git/utils"
)

//var Logger = log.Default()

type registeryServiceImplementation struct {
	UnimplementedRegistryServiceServer
}

func init() {
	log.Printf("registery init")
}

func Start(configData []byte) error {

	lis, err := net.Listen("tcp", ":8502")
	if err != nil {
		log.Fatalf("failed to listen(ServiceBase): %v", err)
		return err
	}

	assignedPort := lis.Addr().(*net.TCPAddr).Port
	log.Printf("assignedPort registry service: %v", assignedPort)

	grpcServer := grpc.NewServer()

	// Register the service implementation with the gRPC server
	registeryService := &registeryServiceImplementation{}
	RegisterRegistryServiceServer(grpcServer, registeryService)
	startListening := func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

	startListening()

	// Return immediately so that the function doesn't block
	return nil
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
