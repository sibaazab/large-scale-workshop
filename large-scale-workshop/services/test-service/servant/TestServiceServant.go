package TestServiceServant

import (
	"math/rand"
	"time"
)

func HelloWorld() string{
    return "Hello World"
}

func HelloToUser(username string) string{
    return "Hello "+ username 
} 

func WaitAndRand(seconds int32, sendToClient func(x int32) error) error {time.Sleep(time.Duration(seconds) * time.Second)
    return sendToClient(int32(rand.Intn(10)))
    }