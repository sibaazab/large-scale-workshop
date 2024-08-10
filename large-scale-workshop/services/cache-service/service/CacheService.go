package CacheService

import (
	"context"

	common "github.com/sibaazab/large-scale-workshop.git/services/common"
	."github.com/sibaazab/large-scale-workshop.git/services/cache-service/common"
	"github.com/sibaazab/large-scale-workshop.git/utils"
	servant "github.com/sibaazab/large-scale-workshop.git/services/cache-service/servant"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type cacheServiceImplementation struct {
	UnimplementedCacheServiceServer
}

func Start(configData []byte) error {
	bindgRPCToService := func(s grpc.ServiceRegistrar) { RegisterCacheServiceServer(s, &cacheServiceImplementation{}) }
	startListening, port, unregister := common.Start("CacheService", 0, bindgRPCToService, messagehandler)
	defer unregister()
	utils.Logger.Printf("CacheService server started on port: %v\n", port)
	servant.CreateChordFromConfig(port)
	startListening()
	return nil
}

func (obj *cacheServiceImplementation) Get(ctxt context.Context, key *wrappers.StringValue) (*wrappers.StringValue, error) {
	val, err := servant.Get(key.Value)
	if err != nil {
		return nil, err
	}
	return wrapperspb.String(val), nil
}

func (obj *cacheServiceImplementation) Set(ctxt context.Context, keyvalue *SetRequest) (*empty.Empty, error) {
	servant.Set(keyvalue.Key, keyvalue.Value)
	return &emptypb.Empty{}, nil
}

func (obj *cacheServiceImplementation) Delete(ctxt context.Context, key *wrappers.StringValue) (*empty.Empty, error) {
	servant.Delete(key.Value)
	return &emptypb.Empty{}, nil
}

func (obj *cacheServiceImplementation) IsAlive(ctxt context.Context, _ *empty.Empty) (*wrappers.BoolValue, error) {
	return wrapperspb.Bool(true), nil
}

func (obj *cacheServiceImplementation) IsRoot(context.Context, *empty.Empty) (*wrappers.BoolValue, error) {
	val, err := servant.IsRoot()
	if err != nil {
		return nil, err
	}
	return wrapperspb.Bool(val), nil
}
