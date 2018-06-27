package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron"
	facebook "github.com/zaingz/sales-and-deals/facebook/post"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	//var wg sync.WaitGroup
	pageIDs := []string{"NishatLinen"}

	cro := cron.New()
	chano := make(chan facebook.Post)
	dano := make(chan bool, 1)

	cro.AddFunc("*/10 * * * *", fetchAllPagePosts(pageIDs, chano, dano))
	cro.Start()
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("test").C("people")

	result := facebook.Posts{}

	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Saved data:", result)

	log.Println("Facebook Post Cron Job Started")
	for {
		select {
		//case _ := <-chano:
		// fmt.Printf("Recived: %v\n\n", post)
		// err = c.Insert(&post)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		case <-dano:
			err = c.Find(bson.M{}).All(&result)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Saved data:", result)
		}
	}

	// wg.Add(1)
	// wg.Wait()

}

func fetchAllPagePosts(pageIDs []string, chnano chan facebook.Post, dano chan bool) func() {
	// config := nsq.NewConfig()
	//p, err := nsq.NewProducer("127.0.0.1:4150", config)
	// if err != nil {
	// 	log.Panic(err)
	// }
	return func() {

		go func() {

			for _, v := range pageIDs {
				var posts *facebook.Posts
				posts = facebook.GetPagePosts(v, 10)

				for _, post := range *posts {
					//payload, err := json.Marshal(post)
					// if err != nil {
					// 	return
					// }

					chnano <- post

					//err = p.Publish("facebook_post", payload)
					// if err != nil {
					// 	log.Panic(err)
					// }
					//log.Println("post published on nsq", post)
				}
				dano <- true

			}
		}()
	}
}
