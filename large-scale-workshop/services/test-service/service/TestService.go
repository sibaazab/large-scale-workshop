package TestService

import (
	"log"
	//"net"
	"context"
	"fmt"

	//"sync"
	//"time"

	//"github.com/pebbe/zmq4"
	services "github.com/sibaazab/large-scale-workshop.git/services/common"
	. "github.com/sibaazab/large-scale-workshop.git/services/test-service/common"

	//. "github.com/sibaazab/large-scale-workshop.git/utils"
	"github.com/sibaazab/large-scale-workshop.git/services/test-service/servant"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
) 
type testServiceImplementation struct{ 
    UnimplementedTestServiceServer 
} 

func Start(configData []byte) error { 
	
    bindgRPCToService := func(s grpc.ServiceRegistrar) { 
        RegisterTestServiceServer(s, &testServiceImplementation{})
    }
    startListening, Port, unregister  :=services.Start("TestService",0, bindgRPCToService, messageHandler) 
	log.Printf("TestService listening on port %d", Port)
	defer unregister()
	startListening()
    return nil
}




func (obj *testServiceImplementation) HelloWorld(ctxt context.Context,_ *emptypb.Empty) (res *wrapperspb.StringValue,err error) {
    return wrapperspb.String(TestServiceServant.HelloWorld()),nil }

func (obj *testServiceImplementation) HelloToUser(ctxt context.Context, in *wrapperspb.StringValue) (res *wrapperspb.StringValue, err error) {
	username := in.GetValue()
	return wrapperspb.String(TestServiceServant.HelloToUser(username)),nil 
}
func (obj *testServiceImplementation) WaitAndRand(seconds *wrapperspb.Int32Value, streamRet TestService_WaitAndRandServer) error {
	log.Printf("WaitAndRand called")
	streamClient := func(x int32) error {
		return streamRet.Send(wrapperspb.Int32(x))
	}
	return TestServiceServant.WaitAndRand(seconds.Value, streamClient)
}

func (obj *testServiceImplementation) Store(ctxt context.Context,kv *StoreKeyValue) (_ *emptypb.Empty,err error) {
	key :=kv.GetKey()
	value := kv.GetValue()
	//servant.cacheMap[key]=vlaue
	err = TestServiceServant.Store(key, value)
	return &emptypb.Empty{}, err
}

func (c *testServiceImplementation) Get(ctxt context.Context,in *wrapperspb.StringValue) (res *wrapperspb.StringValue,err error){
	key := in.GetValue()
	value, err := TestServiceServant.Get(key)
	if err != nil {
		log.Printf("error in test service get")
		return nil, err
	}
	//value:= servant.cacheMap[key]
	return wrapperspb.String(value),nil
}

func (c *testServiceImplementation) IsAlive(ctxt context.Context, _ *emptypb.Empty) (res *wrapperspb.BoolValue,err error){
	return wrapperspb.Bool(TestServiceServant.IsAlive()),nil 
}

func (testServiceImplementation) ExtractLinksFromURL(ctxt context.Context, param *ExtractLinksFromURLParameters) (res *ExtractLinksFromURLReturnedValue,err error) {
	linksArr, err := TestServiceServant.ExtractLinksFromURL(param.GetUrl(), param.GetDepth())
	if err != nil {
		log.Printf("error in test service extract links url")
		return nil, err
	}
	log.Printf("the links array is ")
	return &ExtractLinksFromURLReturnedValue{Links: linksArr}, nil
}





var serviceInstance *testServiceImplementation
 
func handleHelloWorld(ctx context.Context, params []byte) (proto.Message, error) {
	return serviceInstance.HelloWorld(ctx, &emptypb.Empty{})
}

func handleHelloToUser(ctx context.Context, params []byte) (proto.Message, error) {
	in := &wrapperspb.StringValue{}
	if err := proto.Unmarshal(params, in); err != nil {
		return nil, err
	}
	return serviceInstance.HelloToUser(ctx, in)
}

// func handleWaitAndRand(ctx context.Context, params []byte, socket *zmq4.Socket) error {
// 	in := &wrapperspb.Int32Value{}
// 	if err := proto.Unmarshal(params, in); err != nil {
// 		return err
// 	}

// 	// Define a client function to send each result as a separate message
// 	streamClient := func(x int32) error {
// 		// Convert the result to a protobuf message
// 		resultWrapper := &wrapperspb.Int32Value{Value: x}
// 		resultData, err := proto.Marshal(resultWrapper)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal result: %v", err)
// 		}

// 		// Send the result through the ZeroMQ socket
// 		if _, err := socket.SendBytes(resultData, 0); err != nil {
// 			return fmt.Errorf("failed to send result: %v", err)
// 		}

// 		return nil
// 	}

// 	// Call the original WaitAndRand function, passing the custom streamClient
// 	if err := TestServiceServant.WaitAndRand(in.GetValue(), streamClient); err != nil {
// 		return err
// 	}

// 	return nil
// }

func handleStore(ctx context.Context, params []byte) (proto.Message, error) {
	in := &StoreKeyValue{}
	if err := proto.Unmarshal(params, in); err != nil {
		return nil, err
	}
	_, err := serviceInstance.Store(ctx, in)
	return nil, err
}

func handleGet(ctx context.Context, params []byte) (proto.Message, error) {
	in := &wrapperspb.StringValue{}
	if err := proto.Unmarshal(params, in); err != nil {
		return nil, err
	}
	return serviceInstance.Get(ctx, in)
}

func handleIsAlive(ctx context.Context, params []byte) (proto.Message, error) {
	return serviceInstance.IsAlive(ctx, &emptypb.Empty{})
}

func handleExtractLinksFromURL(ctx context.Context, params []byte) (proto.Message, error) {
	in := &ExtractLinksFromURLParameters{}
	if err := proto.Unmarshal(params, in); err != nil {
		return nil, err
	}
	return serviceInstance.ExtractLinksFromURL(ctx, in)
}

func messageHandler(method string, parameters []byte) (response proto.Message, err error) {
	switch method {
	case "HelloWorld":
		return handleHelloWorld(context.Background(), parameters)
	case "HelloToUser":
		return handleHelloToUser(context.Background(), parameters)
	// case "WaitAndRand":
	// 	// Handle streaming case
	// 	if err := handleWaitAndRand(context.Background(), parameters, socket); err != nil {
	// 		return nil, err
	// 	}
	// 	// Return nil because the response is sent within handleWaitAndRand
	// 	return nil, nil
	case "Store":
		return handleStore(context.Background(), parameters)
	case "Get":
		return handleGet(context.Background(), parameters)
	case "IsAlive":
		return handleIsAlive(context.Background(), parameters)
	case "ExtractLinksFromURL":
		return handleExtractLinksFromURL(context.Background(), parameters)
	default:
		return nil, fmt.Errorf("unknown method: %v", method)
	}
}

