package main

import (
	"fmt"
	"sync"

	"github.com/src/saifhatric/chat-v2/handlers"
)

var Wg sync.WaitGroup 
func main(){
  fmt.Println("Server running on port : 8080 ")
  Wg.Add(1)
  func ()  {
    go handlers.ChanListener()
    Wg.Done()
  }()
  
  Wg.Wait()
  handlers.Routes(":8080")
}
