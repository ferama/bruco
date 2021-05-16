package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ferama/coreai/pkg/pool"
)

func main() {
	pool := pool.NewPool(4, "./lambda")
	go func() {
		i := 0
		for {
			i++
			pool.HandleEventAsync([]byte("the event " + fmt.Sprint(i)))
			// res, _ := pool.HandleEvent([]byte("the event " + fmt.Sprint(i)))
			// log.Println(res)
			time.Sleep(time.Second)
		}
	}()

	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	pool.Destroy()
}
