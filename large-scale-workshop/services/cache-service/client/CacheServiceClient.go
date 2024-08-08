 package CacheServiceClient

import (
	"context"
	"fmt"
	//"log"

	//registery  "github.com/sibaazab/large-scale-workshop.git/services/registry-service/service"
	service "github.com/sibaazab/large-scale-workshop.git/services/cache-service/common"
	services "github.com/sibaazab/large-scale-workshop.git/services/common"
	//"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)


type CacheServiceClient struct {
	services.ServiceClientBase[service.CacheServiceClient]
}

func NewCacheServiceClient() *CacheServiceClient {
	client := &CacheServiceClient{
		ServiceClientBase: services.ServiceClientBase[service.CacheServiceClient]{
			ServiceName:  "CacheService",
			CreateClient: service.NewCacheServiceClient,
		},
	}
	client.ServiceClientBase.LoadRegistryAddresses()
	return client
}



func (obj *CacheServiceClient) IsAlive() (bool, error) {
	
	c, closeFunc, _ := obj.Connect()
	defer closeFunc()
	r, err := c.IsAlive(context.Background(), &emptypb.Empty{})

	if err != nil {
		return false, fmt.Errorf("could not call IsAlive: %v", err)
	}

	return r.Value, nil
}



func (obj *CacheServiceClient) Set(key string, value string) error {
	c, closeFunc, err := obj.Connect()

	defer closeFunc()
	_, err = c.Set(context.Background(), &service.SetRequest{Key: key, Value: value})

	if err != nil {
		return fmt.Errorf("could not call Set: %v", err)
	}

	return nil
}

func (obj *CacheServiceClient) Get(key string) (string, error) {
	c, closeFunc, err := obj.Connect()

	defer closeFunc()
	res, err := c.Get(context.Background(), &wrapperspb.StringValue{Value: key})

	if err != nil {
		return "",fmt.Errorf("could not call Get: %v", err)
	}

	return res.Value,nil
}

func (obj *CacheServiceClient) Delete(key string) error {
	c, closeFunc, err := obj.Connect()
	
	defer closeFunc()

	_, err = c.Delete(context.Background(), &wrapperspb.StringValue{Value: key})
	if err != nil {
		return fmt.Errorf("could not call Delete: %v", err)
	}
	return nil
}