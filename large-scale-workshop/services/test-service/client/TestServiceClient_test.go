package TestService

import (
	"context"
	"flag"
	"log"
	"testing"
	//"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/known/wrapperspb"
)


func TestHelloWorld(t *testing.T) {
    c := NewTestServiceClient("localhost:50051") 
    r,err :=c.HelloWorld()
    if err !=nil {
            t.Fatalf("could not call HelloWorld: %v",err)
    return
    }
        t.Logf("Response: %v",r)
}

var username string

var keyGet string
var keyStore string
var value string

// Initialize command-line arguments
func init() {
    flag.StringVar(&username, "username", "", "Username to pass to HelloToUser")
    flag.StringVar(&keyGet, "keyGet", "", "keyGet to pass to get")
    flag.StringVar(&keyStore, "keyStore", "", "key to pass to Store")
    flag.StringVar(&value, "value", "", "value to pass to store")
}

func TestHelloToUser(t *testing.T) {
    flag.Parse()
    if username == "" {
        t.Fatalf("username flag not set")
        return
    }

    c := NewTestServiceClient("localhost:50051")

    res, err := c.HelloToUser(context.Background(), username)
    if err != nil {
        t.Fatalf("could not call HelloToUser: %v", err)
        return
    }
    t.Logf("Response: %v", res)
}

func TestWaitAndRand(t *testing.T) {
    c := NewTestServiceClient("localhost:50051")
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
    

    c := NewTestServiceClient("localhost:50051")

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

    c := NewTestServiceClient("localhost:50051")

    res, err := c.Get(context.Background(), keyGet)
    if err != nil {
        t.Fatalf("could not call Get: %v", err)
        return
    }
    t.Logf("Response: %v", res)
}
