package common

import (
	"fmt"
	"log"
	"net"
	"os"
	//"time"
    "github.com/cenkalti/backoff/v4"

	RegistryServiceClient "github.com/sibaazab/large-scale-workshop.git/services/registry-service/client"
	"gopkg.in/yaml.v2"
	"github.com/pebbe/zmq4"
	"github.com/sibaazab/large-scale-workshop.git/Config"
	. "github.com/sibaazab/large-scale-workshop.git/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)


func startgRPC() (listeningAddress string, grpcServer *grpc.Server, startListening func(), assignedPort int,) {
	lis, err := net.Listen("tcp", ":0")
	assignedPort = lis.Addr().(*net.TCPAddr).Port
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		Logger.Fatalf("failed to listen(ServiceBase): %v", err)
	}

	listeningAddress = lis.Addr().String()
	grpcServer = grpc.NewServer()
	startListening = func() {
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Fatalf("failed to serve: %v", err)
		}
	}
	return
}

func LoadRegistryAddresses() []string {
	configFile := os.Args[1]
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
	
	return config.RegistryAddresses
}


func registerAddress(serviceName string, registryAddresses []string, listeningAddress string) (unregister func()) {
    registryClient := RegistryServiceClient.NewRegistryServiceClient(registryAddresses)
	log.Printf("--------------------------------------------------------------------registerAddress serviceName=%v", serviceName)
	log.Printf("registerAddress registryAddresses=%v", registryAddresses)
	log.Printf("registerAddress listeningAddress=%v", listeningAddress)


	operation := func() error {
        err := registryClient.Register(serviceName, listeningAddress)
        if err != nil {
            log.Printf("Retrying registration for service %v: %v", serviceName, err)
            return err
        }
        return nil
    }

    err := backoff.Retry(operation, backoff.NewExponentialBackOff())
    if err != nil {
        Logger.Fatalf("Failed to register to registry service: %v", err)
    }

    return func() {
        registryClient.Unregister(serviceName, listeningAddress)
    }
}


func StartCache(serviceName string, grpcListenPort int, bindgRPCToService func(grpc.ServiceRegistrar)) (startListening func(), unregister func(), portNum int) {
	listeningAddressNonAsync, grpcServer, startListening, port := startgRPC()
	bindgRPCToService(grpcServer)
	unregisterFunc := registerAddress(serviceName, LoadRegistryAddresses(), listeningAddressNonAsync)

	return startListening, unregisterFunc, port
}

func Start(serviceName string, grpcListenPort int, bindgRPCToService func(grpc.ServiceRegistrar), messageHandler func(method string, parameters []byte) (response proto.Message, err error)) (startListening func(), unregister func(), portNum int, listeningAddressNonAsync string) {
	listeningAddressNonAsync, grpcServer, startListening, port := startgRPC()
	var listeningAddressAsync string
	startMQ, listeningAddressAsync := bindMQToService(0, messageHandler)
	bindgRPCToService(grpcServer)
	unregisterFuncNonAsync := registerAddress(serviceName, LoadRegistryAddresses(), listeningAddressNonAsync)
	unregisterFuncAsync := registerAddress(serviceName+"MQ", LoadRegistryAddresses(), listeningAddressAsync)
	go func() {
		if unregisterFuncAsync != nil {
			defer unregisterFuncAsync()
		}
		if startMQ != nil {
			startMQ()
		}
	}()
	return startListening, unregisterFuncNonAsync, port, listeningAddressNonAsync
}



func bindMQToService(listenPort int, messageHandler func(method string, parameters []byte) (response proto.Message, err error)) (startMQ func(), listeningAddress string){

	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
 		Logger.Fatalf("Failed to create a new zmq socket: %v", err)
	} 
	if listenPort == 0 {
 		listeningAddress = "tcp://127.0.0.1:*"
	} else {
 		listeningAddress = fmt.Sprintf("tcp://127.0.0.1:%v", listenPort)
	} 
	err = socket.Bind(listeningAddress)
	if err != nil {
 		Logger.Fatalf("Failed to bind a zmq socket: %v", err)
	} 
	listeningAddress, err = socket.GetLastEndpoint()
	if err != nil {
		 Logger.Fatalf("Failed to get listetning address of zmq socket: %v", err)
	} 

	startMQ = func() {
		for {
			data, readerr := socket.RecvBytes(0)
			if err != nil {
				Logger.Printf("Failed to receive bytes from MQ socket: %v\n", readerr)
				continue
			}
			if len(data) == 0 {
				continue
			}
			Logger.Printf("data len: %v\n", len(data))


			go func(data []byte) {
				var responseMessage ReturnValue
	
				// Unmarshal the received data into CallParameters
				callParams := &CallParameters{}
				unmarshalErr := proto.Unmarshal(data, callParams)
				if err != nil {
					Logger.Printf("Failed to unmarshal data: %v\n", unmarshalErr)
					responseMessage.Error = unmarshalErr.Error()
					sendResponse(socket, &responseMessage)
					return
				}
	
				// Handle the message using messageHandler
				handlerRes, handlerErr := messageHandler(callParams.Method, callParams.Data)
	
				if handlerErr != nil {
					Logger.Printf("Error in messageHandler: %v\n", handlerErr)
					responseMessage.Error = handlerErr.Error()
				}
	
				// Marshal the response message if available
				if handlerRes != nil {
					marshalData, marshalErr := proto.Marshal(handlerRes)
					if marshalErr != nil {
						Logger.Printf("Failed to marshal response: %v\n", marshalErr)
						responseMessage.Error = marshalErr.Error()
					} else {
						responseMessage.Data = marshalData
					}
				}
	
				// Send the response back to the client
				sendErr := sendResponse(socket, &responseMessage)
				if sendErr != nil {
					Logger.Printf("Failed to send response: %v\n", sendErr)
				}
			}(data)
		}

	}
	return startMQ, listeningAddress
}

// Helper function to send a response through the socket
func sendResponse(socket *zmq4.Socket, responseMessage *ReturnValue) error {
	res, err := proto.Marshal(responseMessage)
	if err != nil {
		Logger.Printf("Failed to marshal responseMessage: %v\n", err)
		return err
	}

	if _, err := socket.SendBytes(res, 0); err != nil {
		Logger.Printf("Failed to send responseMessage: %v\n", err)
		return err
	}

	return nil
}