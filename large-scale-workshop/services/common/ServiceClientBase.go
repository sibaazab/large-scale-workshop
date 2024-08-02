package services

import (
    "context"
    "fmt"
    "math/rand"
    "time"

    "google.golang.org/grpc"
    //"google.golang.org/grpc/credentials/insecure"
    registery_ser "github.com/sibaazab/large-scale-workshop.git/services/registry-service/service"
    "google.golang.org/protobuf/types/known/wrapperspb"
)

type ServiceClientBase[client_t any] struct {
    RegistryAddresses []string
    CreateClient      func(grpc.ClientConnInterface) client_t
    ServiceName string
}

func NewServiceClientBase[client_t any](registryAddresses []string, createClient func(grpc.ClientConnInterface) client_t, serviceName string) *ServiceClientBase[client_t] {
    return &ServiceClientBase[client_t]{
        RegistryAddresses: registryAddresses,
        CreateClient:      createClient,
        ServiceName:      serviceName,
    }
}

func (obj *ServiceClientBase[client_t]) pickNode() (string, error) {

    nodes, err := registery_ser.NewRegistryService().Discover(context.Background(), wrapperspb.String(obj.ServiceName))
    if err != nil {
        return "", err
    }

    rand.Seed(time.Now().UnixNano())
    index := rand.Intn(len(nodes.GetNodes()))
    return nodes.GetNodes()[index], nil
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {
    nodeAddress, err := obj.pickNode()
    if err != nil {
        var empty client_t
        return empty, nil, fmt.Errorf("failed to pick node: %v", err)
    }

    conn, err := grpc.Dial(nodeAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        var empty client_t
        return empty, nil, fmt.Errorf("failed to connect client to %v: %v", nodeAddress, err)
    }
    c := obj.CreateClient(conn)
    return c, func() { conn.Close() }, nil
}