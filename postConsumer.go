package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
	facebook "github.com/zaingz/sales-and-deals/facebook/post"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	decodeConfig := nsq.NewConfig()
	c, err := nsq.NewConsumer("facebook_post", "post-fetch", decodeConfig)
	if err != nil {
		log.Panic("Could not create consumer")
	}
	//c.MaxInFlight defaults to 1

	c.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Println("NSQ message received:")
		var post facebook.Posts
		json.Unmarshal(message.Body, &post)

		log.Println(string(message.Body))
		return nil
	}))

	err = c.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}
	log.Println("Awaiting messages from NSQ")
	wg.Wait()
}
