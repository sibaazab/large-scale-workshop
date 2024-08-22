package RegistryService

import (
	"context"
	"fmt"
	"log"

	//"log"

	registeryClient  "github.com/sibaazab/large-scale-workshop.git/services/commonRegistry"
	service "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
	//"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)


type RegistryServiceClient struct {
	registeryClient.RegistryClientBase
}



func NewRegistryServiceClient(addresses []string) *RegistryServiceClient {
	if len(addresses) <= 0 {
		return nil
	}

	client := &RegistryServiceClient{
		RegistryClientBase: registeryClient.RegistryClientBase{
			Addresses:    addresses,
			CreateClient: service.NewRegistryServiceClient,
		},
	}
	return client
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
	log.Printf("RegisrtyServiceClient %v -------%v", service_name,service_address )
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