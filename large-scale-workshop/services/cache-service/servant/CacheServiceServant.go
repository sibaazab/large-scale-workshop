package CacheServiceServant

import (
	//"context"
	//"time"

	//"log"
	//"net"
	"fmt"
	"os"

	//"sync"
	//"time"
	Config "github.com/sibaazab/large-scale-workshop.git/Config"

	//CacheService "github.com/sibaazab/large-scale-workshop.git/services/cache-service/common"
	common "github.com/sibaazab/large-scale-workshop.git/services/common"
	RegistryServiceClient "github.com/sibaazab/large-scale-workshop.git/services/registry-service/client"
	Chord "github.com/sibaazab/large-scale-workshop.git/services/registry-service/servant/dht"
	"github.com/sibaazab/large-scale-workshop.git/utils"
	"gopkg.in/yaml.v2"
	//"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/known/emptypb"
	//"google.golang.org/protobuf/types/known/wrappers"
)

var chord *Chord.Chord

func CreateChordFromConfig(port int) (*Chord.Chord, error) {
	// Ensure exactly one command-line argument for the configuration file
	if len(os.Args) != 2 {
		return nil, fmt.Errorf("expecting exactly one configuration file")
	}

	// Read configuration file
	fileData, err := os.ReadFile(os.Args[1])
	if err != nil {
		utils.Logger.Fatalf("error reading file: %v", err)
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal YAML data into config object
	var config Config.ConfigBase
	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %w", err)
	}

	// Check for registry addresses
	if len(config.RegistryAddresses) == 0 {
		return nil, fmt.Errorf("no registry addresses found in configuration file")
	}

	// Create RegistryServiceClient
	client := RegistryServiceClient.NewRegistryServiceClient(config.RegistryAddresses)

	// Discover service nodes for "CacheService"
	serviceNodes, err := client.Discover("CacheService")
	if err != nil {
		return nil, fmt.Errorf("error discovering services: %w", err)
	}

	// Initialize or join Chord ring
	if len(serviceNodes) == 0 {
		// No existing service nodes found; create a new Chord ring
		chord, err := Chord.NewChord(fmt.Sprintf(":%v", port), 1099)
		if err != nil {
			return nil, fmt.Errorf("error creating new chord: %w", err)
		}
		return chord, nil
	}

	// Attempt to join an existing Chord ring
	for _, node := range serviceNodes {
		isRoot, err := common.IsRoot(node)
		if err != nil {
			utils.Logger.Fatalf("Error checking if node is root: %v", err)
			return nil, fmt.Errorf("error checking if node is root: %w", err)
		}
		if isRoot {
			chord, err := Chord.JoinChord(fmt.Sprintf(":%v", port), node, 1099)
			if err != nil {
				utils.Logger.Fatalf("Error joining chord: %v", err)
				return nil, fmt.Errorf("error joining chord: %w", err)
			}
			return chord, nil
		}
	}

	// No root node found; this should not happen
	return nil, fmt.Errorf("could not find a root node to join")
}



func Get(key string) (string, error) {
	// Retrieve all keys
	keys, err := chord.GetAllKeys()
	if err != nil {
		utils.Logger.Printf("Failed to get all keys: %v\n", err)
		return "", err
	}

	// Check if the key exists in the retrieved keys
	found := false
	for _, k := range keys {
		if k == key {
			found = true
			break
		}
	}

	if !found {
		return "", nil // Key not found
	}

	// Retrieve the value associated with the key
	res, err := chord.Get(key)
	if err != nil {
		return "", err
	}
	return res, nil
}

func Set(key string, value string) error {
	err := chord.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func Delete(key string) error {
	err := chord.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func IsRoot() (bool, error) {
	val, err := chord.IsFirst()
	if err != nil {
		return false, err
	}
	return val, nil
}
