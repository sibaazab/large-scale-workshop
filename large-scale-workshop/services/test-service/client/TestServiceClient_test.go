package TestService

import (
	"context"
	"flag"
	"log"
	"testing"
	//"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/known/wrapperspb"
)
//var listenAddress string= "[::]:33899"
func TestHelloWorld(t *testing.T) {
	c := NewTestServiceClient()
	r, err := c.HelloWorld()
	if err != nil {
		t.Fatalf("could not call HelloWorld: %v", err)
		return
	}
	t.Logf("Response: %v", r)
}

var username string
var keyGet string
var keyStore string
var value string
var Url string
var Depth32 int32
var Depth int

// Initialize command-line arguments
func init() {
	flag.StringVar(&username, "username", "", "Username to pass to HelloToUser")
	flag.StringVar(&keyGet, "keyGet", "", "keyGet to pass to get")
	flag.StringVar(&keyStore, "keyStore", "", "key to pass to Store")
	flag.StringVar(&value, "value", "", "value to pass to store")
	flag.StringVar(&Url, "Url", "", "Url link")
	flag.IntVar((*int)(&Depth), "Depth", 0, "value to pass to store")
}

func TestHelloToUser(t *testing.T) {
	flag.Parse()
	if username == "" {
		t.Fatalf("username flag not set")
		return
	}

	c := NewTestServiceClient()

	res, err := c.HelloToUser(context.Background(), username)
	if err != nil {
		t.Fatalf("could not call HelloToUser: %v", err)
		return
	}
	t.Logf("Response: %v", res)
}

func TestWaitAndRand(t *testing.T) {
	c := NewTestServiceClient()
	resPromise, err := c.WaitAndRand(3)
	if err != nil {
		t.Fatalf("Calling WaitAndRand failed: %v", err)
		return
	}
	res, err := resPromise()
	if err != nil {
		t.Fatalf("WaitAndRand failed: %v", err)
		return
	}
	t.Logf("Returned random number: %v\n", res)
}

func TestStore(t *testing.T) {
	log.Print("store")
	flag.Parse()
	if keyStore == "" {
		t.Fatalf("key flag not set")
		return
	}
	if value == "" {
		t.Fatalf("value flag not set")
		return
	}

	c := NewTestServiceClient()

	err := c.Store(context.Background(), keyStore, value)
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	t.Logf("Response: success store")
}

func TestGet(t *testing.T) {
	flag.Parse()
	if keyGet == "" {
		t.Fatalf("key flag not set")
		return
	}

	c := NewTestServiceClient()

	res, err := c.Get(context.Background(), keyGet)
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	t.Logf("Response: %v", res)
}

func TestIsAlive(t *testing.T) {
	c := NewTestServiceClient()

	res, err := c.isAlive()
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	t.Logf("Response: %v", res)
}
func TestExtractLinksFromURL(t *testing.T) {
	c := NewTestServiceClient()
	Depth32 = int32(Depth)
	//log.Printf("the depth is %v", Depth32)
	res, err := c.ExtractLinksFromURL(context.Background(), Url, Depth32)
	if err != nil {
		t.Fatalf("could not call ExtractLinksFromURL: client_test (1) %v", err)
		return
	}
	t.Logf("Response: %v", res)
}




func TestHelloWorldAsync(t *testing.T) {
	c := NewTestServiceClient()
	r, err := c.HelloWorldAsync()
	if err != nil {
		t.Fatalf("could not call HelloWorld: %v", err)
		return
	} 
	res, err := r()
	if err != nil {
		t.Fatalf("HelloWorld returned error : %v", err)
		return
	} 
	t.Logf("Response: %v", res)
   } 
   
   func TestHelloToUserAsync(t *testing.T) {
	c := NewTestServiceClient()
	r, err := c.HelloToUserAsync(username)
	if err != nil {
		t.Fatalf("could not start HelloToUserAsync: %v", err)
		return
	}

	// Call the returned function to get the result
	res, err := r()
	if err != nil {
		t.Fatalf("HelloToUserAsync returned error: %v", err)
		return
	}

	// Log the response
	t.Logf("Response: %v", res)
}

func TestWaitAndRandAsync(t *testing.T) {
    c := NewTestServiceClient()
    r, err := c.WaitAndRandAsync(10)
    if err != nil {
        t.Fatalf("could not start WaitAndRandAsync: %v", err)
        return
    }

    // Call the returned function to get the result
    res, err := r()
    if err != nil {
        t.Fatalf("WaitAndRandAsync returned error: %v", err)
        return
    }

    // Log the response
    t.Logf("Response: %v", res)
}

func TestGetAsync(t *testing.T) {
	c := NewTestServiceClient()
	r, err := c.GetAsync(keyGet)
	if err != nil {
		t.Fatalf("could not start GetAsync: %v", err)
		return
	}

	// Call the returned function to get the result
	res, err := r()
	if err != nil {
		t.Fatalf("GetAsync returned error: %v", err)
		return
	}

	// Log the response
	t.Logf("Response: %v", res)
}


func TestStoreAsync(t *testing.T) {
	c := NewTestServiceClient()
	r, err := c.StoreAsync(keyStore, value)
	if err != nil {
		t.Fatalf("could not start StoreAsync: %v", err)
		return
	}

	// Call the returned function to execute the store operation
	err = r()
	if err != nil {
		t.Fatalf("StoreAsync returned error: %v", err)
		return
	}

	// Log success message
	t.Log("Successfully stored the key-value pair")
}


func TestIsAliveAsync(t *testing.T) {
	c := NewTestServiceClient()
	r, err := c.IsAliveAsync()
	if err != nil {
		t.Fatalf("could not start IsAliveAsync: %v", err)
		return
	}

	// Call the returned function to get the result
	res, err := r()
	if err != nil {
		t.Fatalf("IsAliveAsync returned error: %v", err)
		return
	}

	// Log the response
	t.Logf("IsAlive response: %v", res)
}

func TestExtractLinksFromURLAsync(t *testing.T) {
	c := NewTestServiceClient()

	// Use the flag values
	Depth32 = int32(Depth)
	r, err := c.ExtractLinksFromURLAsync(Url, Depth32)
	if err != nil {
		t.Fatalf("could not start ExtractLinksFromURLAsync: %v", err)
		return
	}

	// Call the returned function to get the result
	res, err := r()
	if err != nil {
		t.Fatalf("ExtractLinksFromURLAsync returned error: %v", err)
		return
	}

	// Log the response
	t.Logf("Extracted links: %v", res)
}