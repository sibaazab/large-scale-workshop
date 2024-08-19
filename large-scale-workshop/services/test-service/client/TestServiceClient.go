package TestService

import (
	context "context"
	"fmt"

	//"log"

	//"log"
	services "github.com/sibaazab/large-scale-workshop.git/services/common"
	service "github.com/sibaazab/large-scale-workshop.git/services/test-service/common"

	//"google.golang.org/grpc"
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TestServiceClient struct {
	services.ServiceClientBase[service.TestServiceClient]
}

func NewTestServiceClient() *TestServiceClient {
	client := &TestServiceClient{
		ServiceClientBase: services.ServiceClientBase[service.TestServiceClient]{
			//RegistryAddresses: []string{address},
			CreateClient: service.NewTestServiceClient,
			ServiceName:  "TestService",
		},
	}

	client.ServiceClientBase.LoadRegistryAddresses()
	return client
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
	res, err := client.HelloToUser(ctx, req)
	if err != nil {
		return "", fmt.Errorf("could not call HelloToUser: %v", err)
	}
	return res.GetValue(), nil
}

func (obj *TestServiceClient) WaitAndRand(seconds int32) (func() (int32, error), error) {
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
	res, err := client.Get(ctx, req)
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

func (obj *TestServiceClient) ExtractLinksFromURL(ctx context.Context, url string, depth int32) ([]string, error) {
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer closeFunc()

	req := &service.ExtractLinksFromURLParameters{Url: url, Depth: depth}
	res, err := client.ExtractLinksFromURL(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not call ExtractLinksFromURL: %v", err)
	}
	return res.GetLinks(), nil
}

func (obj *TestServiceClient) HelloWorldAsync() (func() (string, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, err
	}
	msg, err := services.NewMarshaledCallParameter("HelloWorld", &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	_, err = mqsocket.Send(string(msg.Data), zmq.Flag(0))
	if err != nil {
		return nil, err
	}
	// return function (future pattern)
	ret := func() (string, error) {
		defer mqsocket.Close()
		rv, err := mqsocket.Recv(0)
		if err != nil {
			return "", err
		}

		return rv, nil
	}
	return ret, nil
}

func (obj *TestServiceClient) HelloToUserAsync(username string) (func() (string, error), error) {
	// Connect to the service
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}

	// Create the request message
	req := &wrapperspb.StringValue{Value: username}

	// Send the request and get the response asynchronously
	return func() (string, error) {
		defer closeFunc()

		// Define the context for the call
		ctx := context.Background()

		// Call the HelloToUser method
		res, err := client.HelloToUser(ctx, req)
		if err != nil {
			return "", fmt.Errorf("could not call HelloToUser: %v", err)
		}
		return res.GetValue(), nil
	}, nil
}

func (obj *TestServiceClient) WaitAndRandAsync(seconds int32) (func() (int32, error), error) {
	// Connect to the service
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect %v. Error: %v", obj.RegistryAddresses, err)
	}

	// Create the request message with wrapperspb.Int32Value
	req := &wrapperspb.Int32Value{Value: seconds}

	// Return a function that will execute the asynchronous call
	return func() (int32, error) {
		defer closeFunc()

		// Call the WaitAndRand method
		r, err := c.WaitAndRand(context.Background(), req)
		if err != nil {
			return 0, fmt.Errorf("could not call WaitAndRand: %v", err)
		}

		// Receive the result from the stream
		res, err := r.Recv()
		if err != nil {
			return 0, fmt.Errorf("failed to receive result: %v", err)
		}

		return res.GetValue(), nil
	}, nil
}

func (obj *TestServiceClient) GetAsync(key string) (func() (string, error), error) {
	// Connect to the service
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}

	// Create the request message with wrapperspb.StringValue
	req := &wrapperspb.StringValue{Value: key}

	// Return a function that will execute the asynchronous call
	return func() (string, error) {
		defer closeFunc()

		// Call the Get method
		res, err := client.Get(context.Background(), req)
		if err != nil {
			return "", fmt.Errorf("could not call Get: %v", err)
		}

		return res.GetValue(), nil
	}, nil
}

func (obj *TestServiceClient) StoreAsync(key string, value string) (func() error, error) {
	// Connect to the service
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}

	// Create the request message
	req := &service.StoreKeyValue{Key: key, Value: value}

	// Return a function that will execute the asynchronous call
	return func() error {
		defer closeFunc()

		// Call the Store method
		_, err := client.Store(context.Background(), req)
		if err != nil {
			return fmt.Errorf("could not call Store: %v", err)
		}

		return nil
	}, nil
}

func (obj *TestServiceClient) IsAliveAsync() (func() (bool, error), error) {
	// Connect to the service
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}

	// Return a function that will execute the asynchronous call
	return func() (bool, error) {
		defer closeFunc()

		// Call the IsAlive method
		r, err := c.IsAlive(context.Background(), &emptypb.Empty{})
		if err != nil {
			return false, fmt.Errorf("could not call IsAlive: %v", err)
		}

		return r.Value, nil
	}, nil
}

func (obj *TestServiceClient) ExtractLinksFromURLAsync(url string, depth int32) (func() ([]string, error), error) {
	// Connect to the service
	client, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}

	// Create the request message
	req := &service.ExtractLinksFromURLParameters{Url: url, Depth: depth}

	// Return a function that will execute the asynchronous call
	return func() ([]string, error) {
		defer closeFunc()

		// Call the ExtractLinksFromURL method
		res, err := client.ExtractLinksFromURL(context.Background(), req)
		if err != nil {
			return nil, fmt.Errorf("could not call ExtractLinksFromURL: %v", err)
		}

		return res.GetLinks(), nil
	}, nil
}
