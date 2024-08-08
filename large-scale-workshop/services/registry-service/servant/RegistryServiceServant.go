package RegistryServiceServant

import (
	//"context"
	"log"
	"strings"
	//"net"
	"fmt"
	"sync"

	//"time"
	//common "github.com/sibaazab/large-scale-workshop.git/services/common"
	"github.com/sibaazab/large-scale-workshop.git/utils"

	//registeryService "github.com/sibaazab/large-scale-workshop.git/services/registry-service/service"
	//"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/known/emptypb"
	//"google.golang.org/protobuf/types/known/wrapperspb"
	Chord "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant/dht"
)




type RegistryServiceServant struct {
    mu       sync.Mutex
	services map[string]map[string]int
}

var r *RegistryServiceServant

var chord *Chord.Chord

func init() {
	r = &RegistryServiceServant{services: make(map[string]map[string]int)}
	log.Printf("map initialized- registery service")
}

func CreateChord(basePort, listenPort int) error {
	var err   error

	address := fmt.Sprintf(":%v", listenPort)
	baseAddress := fmt.Sprintf(":%v", basePort)

	if listenPort == basePort {
		// Create a new Chord ring
		chord, err = Chord.NewChord(address, 1099)
		if err != nil {
			return fmt.Errorf("failed to create Chord ring: %w", err)
		}
		utils.Logger.Printf("Chord ring created on %s", address)
	} else {
		// Join an existing Chord ring
		chord, err = Chord.JoinChord(address, baseAddress, 1099)
		if err != nil {
			return fmt.Errorf("failed to join Chord ring at %s: %w", baseAddress, err)
		}
		utils.Logger.Printf("Joined Chord ring on %s via %s", address, baseAddress)
	}

	// If needed, you could return the chord instance, or handle it differently
	// return chord, nil
	return nil
}




func Register(serviceName string, address string) error {
	
	// Get the current addresses associated with the service
	res, err := chord.Get(serviceName)
	if err != nil {
		return err
	}

	// Decode the result into a slice of addresses
	addresses := strings.Split(res, ",")
	if res == "" {
		addresses = []string{}
	}

	// Append the new address and update the DHT entry
	addresses = append(addresses, address)
	err = chord.Set(serviceName, strings.Join(addresses, ","))
	if err != nil {
		utils.Logger.Printf("Failed to register service %v on address %v\n", serviceName, address)
		return err
	}

	return nil
}



func Unregister(serviceName string, address string) error {
	res, err := chord.Get(serviceName)
	if err != nil {
		return err
	}

	// Decode the result directly
	addresses := strings.Split(res, ",")
	if len(addresses) == 0 {
		return nil
	}

	// Find and remove the address
	for i, addr := range addresses {
		if addr == address {
			addresses = append(addresses[:i], addresses[i+1:]...)
			break
		}
	}

	// Update or delete the entry in the DHT
	if len(addresses) == 0 {
		err = chord.Delete(serviceName)
	} else {
		err = chord.Set(serviceName, strings.Join(addresses, ","))
	}

	return err
}


func Discover(serviceName string) []string {
	res, err := chord.Get(serviceName)
	if err != nil {
		utils.Logger.Printf("Failed to discover services for %s: %v", serviceName, err)
		return []string{}
	}
    
	if res == "" {
		return []string{}
	}
	addresses := strings.Split(res, ",")

	return addresses
}



