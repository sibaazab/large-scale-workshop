package commonRegistry


import(
	"fmt"
	"log"

	//"log"
	"math/rand"
	service "github.com/sibaazab/large-scale-workshop.git/services/registry-service/common"
	"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/known/wrapperspb"
)

type RegistryClientBase struct {
	Addresses    []string
	CreateClient func(grpc.ClientConnInterface) service.RegistryServiceClient
} 

func (obj *RegistryClientBase) PickRandomRegistry() string {
	// Pick a random index
	index := rand.Intn(len(obj.Addresses))

	//log.Printf("picked the registery, %v", obj.Addresses[index])
	log.Printf("the registery chosen forn the testArevice is %v", obj.Addresses)
	return obj.Addresses[index]
}

func (obj *RegistryClientBase) Connect() (res service.RegistryServiceClient, closeFunc func(), err error) {
	regAddress := obj.PickRandomRegistry()
	conn, err := grpc.Dial(regAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty service.RegistryServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", regAddress, err)
	}
	c := obj.CreateClient(conn)
	log.Printf("the registry was picked is %v, Connect()- ResistryServiceClient", regAddress)
	return c, func() { conn.Close() }, nil
}