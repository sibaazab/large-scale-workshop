package common

import (
	//"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	//"google.golang.org/grpc/credentials/insecure"
	registeryClient "github.com/sibaazab/large-scale-workshop.git/services/registry-service/client"
	//"google.golang.org/protobuf/types/known/wrapperspb"
    "github.com/sibaazab/large-scale-workshop.git/Config"
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

func (obj *ServiceClientBase[client_t]) LoadRegistryAddresses() {
	configFile := "/workspaces/large-scale-workshop/services/common/RegistryAddresses.yaml"
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading registry yaml file: %v", err)
		os.Exit(2)
	}
	var config Config.RegistryServiceConfig
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling registry addresses data: %v", err)
		os.Exit(3)
	}
	if len(config.RegistryAddresses) <= 0 {
		log.Fatalf("registry addresses yaml file does not include any enteries")
		os.Exit(4)
	}
	obj.RegistryAddresses = config.RegistryAddresses
}

func (obj *ServiceClientBase[client_t]) pickNode() (string, error) {

    //nodes, err := registery_ser.Discover(context.Background(), wrapperspb.String(obj.ServiceName))
	client:= registeryClient.NewRegistryServiceClient(obj.RegistryAddresses)
    nodes, err := client.Discover(obj.ServiceName)
	

	if err != nil {
        return "", err
    }

    rand.Seed(time.Now().UnixNano())
    index := rand.Intn(len(nodes))
	log.Printf("%v", nodes[index])
    return nodes[index], nil
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