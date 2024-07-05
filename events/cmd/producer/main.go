package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/drawiin/go-expert/events/pkg/rabbitmq"
)

func main() {
	ch := rabbitmq.OpenChannel()
	defer ch.Close()

	fmt.Println("Starting to send messages...")
	var i int
	for {
		sleep :=  rand.Intn(5)
		time.Sleep(time.Duration(sleep) * time.Second)
		i++
		times := rand.Intn(100) + 1
		fmt.Print("Sending ", times, " messages... ")
		for j := 0; j < times; j++ {
			body := fmt.Sprintf("Message %d-%d", i, j+1)
			err := rabbitmq.Publish(ch, body, "amq.direct")
			if err != nil {
				fmt.Print(err)
			} else {
				fmt.Print(body, ", ")
			}
		}
		fmt.Println("=============================================")

	}
}
