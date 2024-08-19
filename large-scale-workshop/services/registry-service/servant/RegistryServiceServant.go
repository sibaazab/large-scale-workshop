package RegistryServiceServant

import (
	//"context"
	"log"
	"strings"
	"time"

	//"net"
	"fmt"
	"sync"

	//"time"
	common "github.com/sibaazab/large-scale-workshop.git/services/common"
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
var isAliveCheck map[string]int
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




func CheckIsAlive() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			utils.Logger.Printf("Started IsAlive check\n")
			CheckAllNodesStatus()
		}()
	}
}

func CheckIfKeyInKeys(key string) {
	keys, err := chord.GetAllKeys()
	if err != nil {
		utils.Logger.Printf("Failed to get all keys: %v\n", err)
		return
	}
	if contains(keys, key) {
		return
	}
	if err := chord.Set(key, ""); err != nil {
		utils.Logger.Printf("Failed to set key %v: %v\n", key, err)
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func DeleteByValue(slice []string, x ...string) []string {
	var result []string
	for _, v := range slice {
		delete := false
		for _, del := range x {
			if v == del {
				delete = true
				break
			}
		}
		if !delete {
			result = append(result, v)
		}
	}
	return result
}


func CheckAllNodesStatus() {
	services, err := chord.GetAllKeys()
	if err != nil {
		utils.Logger.Printf("Failed to get all keys: %v\n", err)
		return
	}

	for _, service := range services {
		nodes, err := chord.Get(service)
		if err != nil {
			utils.Logger.Printf("Failed to get nodes for service %v: %v\n", service, err)
			continue
		}

		var nodesArr []string
		if nodes != "" {
			nodesArr = strings.Split(nodes, ",")
		}

		if len(nodesArr) == 0 {
			continue
		}

		nodesToDelete := []string{}

		for _, nodeAddr := range nodesArr {
			client := common.NewCacheServiceBase(nodeAddr)
			prevFailures := isAliveCheck[nodeAddr]
			isAlive, err := client.IsAlive(service)
			if err != nil || !isAlive {
				utils.Logger.Printf("Node %v of service %v is not alive.\n", nodeAddr, service)
				prevFailures++
			} else {
				prevFailures = 0
			}

			isAliveCheck[nodeAddr] = prevFailures

			if prevFailures >= 2 {
				nodesToDelete = append(nodesToDelete, nodeAddr)
			}
		}

		if len(nodesToDelete) > 0 {
			updateNodes(service, nodesArr, nodesToDelete)
		}
	}
}

func updateNodes(service string, nodesArr []string, nodesToDelete []string) {
	nodesArr = DeleteByValue(nodesArr, nodesToDelete...)
	if len(nodesArr) == 0 {
		if err := chord.Delete(service); err != nil {
			utils.Logger.Printf("Failed to delete service %v: %v\n", service, err)
		}
	} else {
		nodesStr := strings.Join(nodesArr, ",")
		if err := chord.Set(service, nodesStr); err != nil {
			utils.Logger.Printf("Failed to update service %v with nodes %v: %v\n", service, nodesArr, err)
		}
	}
	utils.Logger.Printf("Service %v has the following nodes: %v\n", service, nodesArr)
}
