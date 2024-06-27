package TestService

import (
	"testing"
    "context"
    "flag"
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

// Initialize command-line arguments
func init() {
    flag.StringVar(&username, "username", "", "Username to pass to HelloToUser")
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