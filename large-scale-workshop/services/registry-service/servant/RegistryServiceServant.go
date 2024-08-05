package RegistryServiceServant


import (
    //"context"
    "log"
    //"net"
    //"fmt"
    "sync"
    //"time"
    //common "github.com/sibaazab/large-scale-workshop.git/services/common"

    //registeryService "github.com/sibaazab/large-scale-workshop.git/services/registry-service/service"
    //"google.golang.org/grpc"
    //"google.golang.org/protobuf/types/known/emptypb"
    //"google.golang.org/protobuf/types/known/wrapperspb"
)



type RegistryServiceServant struct {
    mu       sync.Mutex
	services map[string]map[string]int
}

var r *RegistryServiceServant

func init() {
	r = &RegistryServiceServant{services: make(map[string]map[string]int)}
	log.Printf("map initialized- registery service")
}

func Register( serviceName string, address string) error {
	r.mu.Lock()
    defer r.mu.Unlock()
    log.Printf("RegisteryService")
    if _, exists := r.services[serviceName]; !exists {
        r.services[serviceName] = make(map[string]int)
    }
    r.services[serviceName][address] = 0

    return nil
}


func Unregister( serviceName string, address string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if nodes, exists := r.services[serviceName]; exists {
        delete(nodes, address)
        if len(nodes) == 0 {
            delete(r.services, serviceName)
        }
    }
    return nil
}

func Discover(serviceName string)  []string {
    r.mu.Lock()
    defer r.mu.Unlock()

    nodes := []string{}
    if serviceNodes, exists := r.services[serviceName]; exists {
        for nodeAddress := range serviceNodes {
            nodes = append(nodes, nodeAddress)
        }
    }

    return nodes
}