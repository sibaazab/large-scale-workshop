package TestService

import (
	context "context"
	"fmt"

	//"log"
	services "github.com/sibaazab/large-scale-workshop.git/services/common"
	service "github.com/sibaazab/large-scale-workshop.git/services/test-service/common"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TestServiceClient struct {
	services.ServiceClientBase[service.TestServiceClient]
}

func NewTestServiceClient(address string) *TestServiceClient {
	return &TestServiceClient{
		ServiceClientBase: services.ServiceClientBase[service.TestServiceClient]{
			RegistryAddresses: []string{address},
			CreateClient: func(conn grpc.ClientConnInterface) service.TestServiceClient {
				return service.NewTestServiceClient(conn)
			},
			ServiceName: "TestService",
		},
	}
}
func (obj *TestServiceClient) HelloWorld() (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect: %v", err)
	}
	
	defer closeFunc()
	// Call the HelloWorld RPC function
	r, err := c.HelloWorld(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("could not call HelloWorld: %v", err)
	}
	return r.Value, nil
}

func (obj *TestServiceClient) HelloToUser(ctx context.Context, username string) (string, error) {
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect: %v", err)
	}
	defer closeFunc()

	req := &wrapperspb.StringValue{Value: username}
	res, err := client.HelloToUser(ctx ,req)
	if err != nil {
		return "", fmt.Errorf("could not call HelloToUser: %v", err)
	}
	return res.GetValue(), nil
}
 
func (obj *TestServiceClient) WaitAndRand(seconds int32) (func() (int32,error), error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect %v. Error: %v", obj.RegistryAddresses, err)
	}
	r, err := c.WaitAndRand(context.Background(), wrapperspb.Int32(seconds))
	if err != nil {
		return nil, fmt.Errorf("could not call Get: %v", err)
	}
	res := func() (int32, error) {
	defer closeFunc()
	x, err := r.Recv()
		return x.Value, err
	}
	return res, nil
}

func (obj *TestServiceClient) Get(ctx context.Context, key string) (string, error) {
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect: %v", err)
	}
	defer closeFunc()

	req := &wrapperspb.StringValue{Value: key}
	res, err := client.Get(ctx ,req)
	if err != nil {
		return "", fmt.Errorf("could not call HelloToUser: %v", err)
	}
	return res.GetValue(), nil
}
 
func (obj *TestServiceClient) Store(ctx context.Context, key string, value string) error {
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}
	defer closeFunc()

	req := &service.StoreKeyValue{Key: key, Value: value}
	_, err = client.Store(ctx, req)
	if err != nil {
		return fmt.Errorf("could not call Store: %v", err)
	}
	return nil
}


func (obj *TestServiceClient) isAlive() (bool, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return false, fmt.Errorf("could not connect: %v", err)
	}
	defer closeFunc()

	r, err := c.IsAlive(context.Background(), &emptypb.Empty{})
	if err != nil {
		return false, fmt.Errorf("could not call isAlive: %v", err)
	}
	return r.Value, nil
}

func (obj *TestServiceClient) ExtractLinksFromURL(ctx context.Context, url string, depth int32) ([]string,error) {
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return nil,fmt.Errorf("could not connect: %v", err)
	}
	defer closeFunc()
 
	req := &service.ExtractLinksFromURLParameters{Url: url, Depth: depth}
	res, err:= client.ExtractLinksFromURL(ctx ,req)
	if err != nil {
		return nil, fmt.Errorf("could not call ExtractLinksFromURL: %v", err)
	}
	return res.GetLinks(), nil
}

