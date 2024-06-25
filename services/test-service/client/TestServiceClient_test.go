package TestService
import(

"testing"

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