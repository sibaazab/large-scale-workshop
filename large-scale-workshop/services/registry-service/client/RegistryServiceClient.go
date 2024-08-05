package ResgistryService

import (
	"context"
	"fmt"
	"log"

	//"log"

	//"log"

	"math/rand"
	//registery  "github.com/sibaazab/large-scale-workshop.git/services/registry-service/service"
	service "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)


type RegistryServiceClient struct {
	Addresses    []string
	CreateClient func(grpc.ClientConnInterface) service.RegistryServiceClient
}

func (obj *RegistryServiceClient) PickRandomRegistry() string {
	// Pick a random index
	index := rand.Intn(len(obj.Addresses))

	// Get the random element
	//log.Printf("picked the registery, %v", obj.Addresses[index])
	log.Printf("the registery chosen forn the testArevice is %v", obj.Addresses[index])
	return obj.Addresses[index]
}

func (obj *RegistryServiceClient) Connect() (res service.RegistryServiceClient, closeFunc func(), err error) {
	regAddress := obj.PickRandomRegistry()
	conn, err := grpc.Dial(regAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty service.RegistryServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", regAddress, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}

func NewRegistryServiceClient(addresses []string) *RegistryServiceClient {
	if len(addresses) <= 0 {
		return nil
	}

	return &RegistryServiceClient{
		Addresses:    addresses,
		CreateClient: service.NewRegistryServiceClient,
	}
}

func (obj *RegistryServiceClient) Discover(service_name string) ([]string, error) {
	
	c, closeFunc, _ := obj.Connect()
	defer closeFunc()
	r, err := c.Discover(context.Background(), &wrapperspb.StringValue{Value: service_name})

	if err != nil {
		return nil, fmt.Errorf("could not call Discover: %v", err)
	}

	return r.Nodes, nil
}

func (obj *RegistryServiceClient) Register(service_name string, service_address string) error {
	c, closeFunc, _ := obj.Connect()

	defer closeFunc()
	log.Printf("RegisrtyServiceClient")
	_, err := c.Register(context.Background(), &service.ServiceRequest{Name: service_name, Address: service_address})

	if err != nil {
		return fmt.Errorf("could not call Register: %v", err)
	}

	return nil
}

func (obj *RegistryServiceClient) Unregister(service_name string, service_address string) error {
	c, closeFunc, _ := obj.Connect()
	defer closeFunc()

	_, err := c.Unregister(context.Background(), &service.ServiceRequest{Name: service_name, Address: service_address})

	if err != nil {
		return fmt.Errorf("could not call Discover: %v", err)
	}

	return nil
}