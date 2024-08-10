package common

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// IsAlive checks if the service is alive by invoking the IsAlive gRPC method.
func IsAlive(address, serviceName string) (bool, error) {
	fullMethodName := fmt.Sprintf("/%s.%s/IsAlive", strings.ToLower(serviceName), serviceName)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*5))
	if err != nil {
		return false, fmt.Errorf("failed to connect to %v: %v", address, err)
	}
	defer conn.Close()

	client := new(wrappers.BoolValue)
	err = conn.Invoke(context.Background(), fullMethodName, &emptypb.Empty{}, client)
	if err != nil {
		return false, err
	}
	return client.Value, nil
}

// IsRoot checks if the node is the root node by invoking the IsRoot gRPC method.
func IsRoot(address string) (bool, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*10))
	if err != nil {
		return false, fmt.Errorf("failed to connect to %v: %v", address, err)
	}
	defer conn.Close()

	client := new(wrappers.BoolValue)
	err = conn.Invoke(context.Background(), "/cacheservice.CacheService/IsRoot", &emptypb.Empty{}, client)
	if err != nil {
		return false, err
	}
	return client.Value, nil
}