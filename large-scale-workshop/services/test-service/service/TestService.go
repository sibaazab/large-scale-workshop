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
var cacheMap map[string]string

func init() {
		cacheMap = make(map[string]string)
		log.Printf("map initialized")
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
	vlaue := kv.GetValue()
	cacheMap[key]=vlaue
	return
}

func (c *testServiceImplementation) Get(ctxt context.Context,in *wrapperspb.StringValue) (res *wrapperspb.StringValue,err error){
	key := in.GetValue()
	value:= cacheMap[key]
	return wrapperspb.String(TestServiceServant.Get(value)),nil
}

func (c *testServiceImplementation) IsAlive(ctxt context.Context, _ *emptypb.Empty) (res *wrapperspb.BoolValue,err error){
	return wrapperspb.Bool(TestServiceServant.IsAlive()),nil 
}
