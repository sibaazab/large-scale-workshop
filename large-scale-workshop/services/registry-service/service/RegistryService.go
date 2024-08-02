package ResgistryService

import (
    "context"
    "log"
    "net"
    //"fmt"
    "sync"
    "time"
    common "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
    "google.golang.org/grpc"
    "google.golang.org/protobuf/types/known/emptypb"
    "google.golang.org/protobuf/types/known/wrapperspb"
)

type RegistryService struct {
    mu       sync.Mutex
    services map[string]map[string]int
    common.UnimplementedRegistryServiceServer
}

func NewRegistryService() *RegistryService {
    return &RegistryService{
        services: make(map[string]map[string]int),
    }
}
func init(){
    log.Printf("registery init")
}

func Start(configData []byte) error {

    //bindgRPCToService := func(s grpc.ServiceRegistrar) { 
     //   RegisterTestServiceServer(s, &testServiceImplementation{})
    
    lis, err := net.Listen("tcp", ":8502")
    assignedPort := lis.Addr().(*net.TCPAddr).Port
    log.Printf("assignedPort registery service: %v", assignedPort)
     if err != nil {
         log.Fatalf("failed to listen(ServiceBase): %v", err)
     }
 
    //listeningAddress := lis.Addr().String()
    grpcServer := grpc.NewServer()
    startListening := func() {
        if err := grpcServer.Serve(lis); err != nil {
             log.Fatalf("failed to serve: %v", err)
         }
    }
    startListening()

	
    return nil
}

func (r *RegistryService) Register(ctx context.Context, req *common.ServiceRequest) (*emptypb.Empty, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.services[req.Name]; !exists {
        r.services[req.Name] = make(map[string]int)
    }
    r.services[req.Name][req.Address] = 0

    go r.monitorNode(req.Name, req.Address, ctx)
    return &emptypb.Empty{}, nil
}

func (r *RegistryService) Unregister(ctx context.Context, req *common.ServiceRequest) (*emptypb.Empty, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if nodes, exists := r.services[req.Name]; exists {
        delete(nodes, req.Address)
        if len(nodes) == 0 {
            delete(r.services, req.Name)
        }
    }
    return &emptypb.Empty{}, nil
}

func (r *RegistryService) Discover(ctx context.Context, req *wrapperspb.StringValue) (*common.ServiceNodes, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    nodes := []string{}
    if serviceNodes, exists := r.services[req.GetValue()]; exists {
        for nodeAddress := range serviceNodes {
            nodes = append(nodes, nodeAddress)
        }
    }

    return &common.ServiceNodes{Nodes: nodes}, nil
}

func (r *RegistryService) monitorNode(serviceName, nodeAddress string, ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        // Assuming the IsAlive method is a placeholder for actual health check logic
        isAlive := r.IsAlive(ctx)
        if !isAlive {
            r.mu.Lock()
            r.services[serviceName][nodeAddress]++
            if r.services[serviceName][nodeAddress] >= 2 {
                log.Printf("Unregistering node %s for service %s due to consecutive failures", nodeAddress, serviceName)
                delete(r.services[serviceName], nodeAddress)
                if len(r.services[serviceName]) == 0 {
                    delete(r.services, serviceName)
                }
                r.mu.Unlock()
                return
            }
            r.mu.Unlock()
        } else {
            r.mu.Lock()
            r.services[serviceName][nodeAddress] = 0
            r.mu.Unlock()
        }
    }
}

func (r *RegistryService) IsAlive(ctx context.Context) bool {
    return true
}

