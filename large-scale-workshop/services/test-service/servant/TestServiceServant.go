package TestServiceServant

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	metaffi "github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"github.com/sibaazab/large-scale-workshop.git/utils"
	CacheServiceClient "github.com/sibaazab/large-scale-workshop.git/services/cache-service/client"
)

var pythonRuntime *metaffi.MetaFFIRuntime
var crawlerModule *metaffi.MetaFFIModule
var extract_links_from_url func(...interface{}) ([]interface{}, error)

var cacheMap map[string]string

func init() {
	cacheMap = make(map[string]string)
	log.Printf("map initialized sdfs")

	pythonRuntime = metaffi.NewMetaFFIRuntime("python311")
	err := pythonRuntime.LoadRuntimePlugin()
	if err != nil {
		msg := fmt.Sprintf("Failed to load runtime plugin: %v", err)
		utils.Logger.Fatalf(msg)
		panic(msg)
	}
	// Load the Crawler module
	crawlerModule, err = pythonRuntime.LoadModule("/workspaces/large-scale-workshop/services/test-service/servant/crawler.py")
	if err != nil {
		msg := fmt.Sprintf("Failed to load ./crawler/crawler.py module: %v", err)
		utils.Logger.Fatalf(msg)
		panic(msg)
	}
	// Load the crawler function
	extract_links_from_url, err = crawlerModule.Load("callable=extract_links_from_url",
		[]IDL.MetaFFIType{IDL.STRING8, IDL.INT64}, // parameters types
		[]IDL.MetaFFIType{IDL.STRING8_ARRAY})      // return type
	if err != nil {
		msg := fmt.Sprintf("Failed to load extract_links_from_url function: %v", err)
		utils.Logger.Fatalf(msg)
		panic(msg)
	}
}

func HelloWorld() string {
	return "Hello World"
}

func HelloToUser(username string) string {
	return "Hello " + username
}

func WaitAndRand(seconds int32, sendToClient func(x int32) error) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return sendToClient(int32(rand.Intn(10)))
}

func Get(key string) (string, error) {
	c := CacheServiceClient.NewCacheServiceClient()
	value, err := c.Get(key)
	if err != nil {
		return "", fmt.Errorf("Failed to get from cache: %v", err)
	}
	return value, nil
}

func Store(key string, value string) error {
	c := CacheServiceClient.NewCacheServiceClient()
	err := c.Set(key, value)
	if err != nil {
		return fmt.Errorf("Failed to store in cache: %v", err)
	}
	return nil
}

func IsAlive() bool {
	return true
}

func ExtractLinksFromURL(url string, depth int32) ([]string, error) {
	// Call Python's extract_links_from_url.
	res, err := extract_links_from_url(url, int64(depth))
	if err != nil {
		log.Printf("error in servant, python returns error")
		return nil, err
	}
	if depth == 0 {
		//empty_lst:= make([]string, 0,0)

		lst := make([]string, 0, 1)
		lst = append(lst, url)
		return lst, nil
	}
	return res[0].([]string), nil
}
