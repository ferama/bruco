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
	workers := pool.NewPool(8, "./lambda")
	go func() {
		i := 0
		for {
			i++
			callback := func(msg *pool.Response) {
				log.Println(msg.Data)
			}
			workers.HandleEvent([]byte("the event "+fmt.Sprint(i)), callback)
			if i%8 == 0 {
				time.Sleep(8 * time.Second)
			}
			// time.Sleep(time.Second / 2)
		}
	}()

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	workers.Destroy()
}
