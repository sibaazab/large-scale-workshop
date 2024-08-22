package common

import (
	//"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	//"time"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/pebbe/zmq4"
	registryClient "github.com/sibaazab/large-scale-workshop.git/services/registry-service/client"
	"github.com/sibaazab/large-scale-workshop.git/Config"
	"google.golang.org/protobuf/proto"
)

type ServiceClientBase[client_t any] struct {
    RegistryAddresses  []string
    CreateClient      func(grpc.ClientConnInterface) client_t
    ServiceName string
}

func (obj *ServiceClientBase[client_t]) LoadRegistryAddresses() {
	configFile := "github.com/sibaazab/large-scale-workshop.git/services/common/RegistryAddresses.yaml"
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

func pickNode(nodes []string) (string, error) {

	log.Printf("nodes= %v", nodes)
    index := rand.Intn(len(nodes))
	log.Printf("%v", nodes[index])
    return nodes[index], nil
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {

	client:= registryClient.NewRegistryServiceClient(obj.RegistryAddresses)
    nodes, err := client.Discover(obj.ServiceName)
	log.Printf("5555555555555555555555555555555555555555555555555")
	

	if err != nil {
		var empty client_t
        return empty,nil, err
    }

    nodeAddress, err := pickNode(nodes)
	log.Printf("5555555555555555555555555555555555555555555555555")

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



func (obj *ServiceClientBase[client_t]) getMQNodes() ([]string, error) {
	registryClient := registryClient.NewRegistryServiceClient([]string{"registryAddress"})
	nodes, err := registryClient.Discover(obj.ServiceName + "MQ")
	if err != nil {
		return nil, fmt.Errorf("failed to discover MQ nodes: %v", err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no MQ nodes found for service: %v", obj.ServiceName)
	}

	return nodes, nil
}



func (obj *ServiceClientBase[client_t]) ConnectMQ() (socket *zmq4.Socket, err error) {
	// Create a new ZeroMQ socket
	socket, err = zmq4.NewSocket(zmq4.REQ) // REQ socket type for sending requests and receiving replies
	if err != nil {
		log.Fatalf("Failed to create a new zmq socket: %v", err)
		return nil, err
	}

	// Get the list of MQ nodes to connect to
	nodes, err := obj.getMQNodes()
	if err != nil {
		log.Fatalf("Failed to get MQ nodes: %v", err)
		return nil, err
	}

	// Connect to each node
	for _, node := range nodes {
		err = socket.Connect(node)
		if err != nil {
			log.Printf("Failed to connect to node %v: %v\n", node, err)
		}
	}

	return socket, nil
}


// MarshaledCallParameter is a structure representing the marshaled call parameters.
type MarshaledCallParameter struct {
    Method string
    Data   []byte
}

// NewMarshaledCallParameter creates a new MarshaledCallParameter by serializing the given method and message.
func NewMarshaledCallParameter(method string, msg proto.Message) (*MarshaledCallParameter, error) {
    // Check if method is valid
    if method == "" {
        return nil, errors.New("method name cannot be empty")
    }
    
    // Serialize the message to a byte slice
    data, err := proto.Marshal(msg)
    if err != nil {
        return nil, err
    }
    
    // Create and return the MarshaledCallParameter
    return &MarshaledCallParameter{
        Method: method,
        Data:   data,
    }, nil
}