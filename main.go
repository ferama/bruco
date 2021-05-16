package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ferama/coreai/pkg/pool"
)

func main() {
	workers := pool.NewPool(4, "./lambda")
	go func() {
		i := 0
		for {
			i++
			callback := func(msg *pool.Response) {
				log.Println(msg.Data)
			}
			workers.HandleEvent([]byte("the event "+fmt.Sprint(i)), callback)
			time.Sleep(time.Second)
		}
	}()

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	workers.Destroy()
}
