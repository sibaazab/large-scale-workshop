package TestService 
import( 
    "log" 
    "net" 
    services "github.com/sibaazab/large-scale-workshop.git/services/common" 
    . "github.com/sibaazab/large-scale-workshop.git/services/test-service/common" 
    "google.golang.org/grpc" 
) 
type testServiceImplementation struct{ 
    UnimplementedTestServiceServer 
} 
func Start(configData []byte) error { 
bindgRPCToService := func(s grpc.ServiceRegistrar) { 
RegisterTestServiceServer(s, &testServiceImplementation{}) 
} 
services.Start("TestService", 50051, bindgRPCToService) 
}