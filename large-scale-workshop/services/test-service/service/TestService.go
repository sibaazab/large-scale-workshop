package TestService

import (
	"log"
	//"net"
	"context"
	
	services "github.com/sibaazab/large-scale-workshop.git/services/common"
	. "github.com/sibaazab/large-scale-workshop.git/services/test-service/common"

	//. "github.com/sibaazab/large-scale-workshop.git/utils"
	"github.com/sibaazab/large-scale-workshop.git/services/test-service/servant"
	"google.golang.org/grpc"
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
    services.Start("TestService", 50051, bindgRPCToService) 
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
	TestServiceServant.Store(key, value)
	return
}

func (c *testServiceImplementation) Get(ctxt context.Context,in *wrapperspb.StringValue) (res *wrapperspb.StringValue,err error){
	key := in.GetValue()
	//value:= servant.cacheMap[key]
	return wrapperspb.String(TestServiceServant.Get(key)),nil
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