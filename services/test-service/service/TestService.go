package TestService

import (
	//"log"
	//"net"
	"context"
    
	services "github.com/sibaazab/large-scale-workshop.git/services/common"
	. "github.com/sibaazab/large-scale-workshop.git/services/test-service/common"
	//. "github.com/sibaazab/large-scale-workshop.git/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
    "github.com/sibaazab/large-scale-workshop.git/services/test-service/servant"
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