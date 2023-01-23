package main

import (
	"log"
	"time"
)

func main() {
	<-Or(Sig(2*time.Second), Sig(5*time.Second), Sig(10*time.Second), Sig(1*time.Minute))
	log.Println("The end")
}
func Sig(period time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(period)
	}()
	return c
}
func Or(ch ...<-chan interface{}) <-chan interface{} {
	recv := make(chan interface{})
	for i, v := range ch {
		go func(channel <-chan interface{}, ind int) {
			recv <- channel
			log.Print("index of chan: ", ind)
		}(v, i)
	}
	return recv
}
