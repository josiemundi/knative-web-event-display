package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	eventsource "gopkg.in/antage/eventsource.v1"

	"fmt"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	"knative.dev/eventing-contrib/pkg/kncloudevents"
)

var es eventsource.EventSource
var id int

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func display(event cloudevents.Event) {
	fmt.Printf("☁️  cloudevents.Event\n%s", event.String())
	fmt.Println(es)
	es.SendEventMessage(event.String(), "tick-event", strconv.Itoa(id))
	id++
}

func main() {
	es = eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/events", es)
	http.HandleFunc("/test", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	// go func() {
	// 	id := 1
	// 	for {
	// 		es.SendEventMessage("tick", "tick-event", strconv.Itoa(id))
	// 		id++
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}

	fmt.Println("Starting")
	go func() {
		fmt.Println("Starting to Listen and Serve")
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}()
	time.Sleep(4000 * time.Millisecond)
	fmt.Println("Middle")
	log.Fatal(c.StartReceiver(context.Background(), display))
	fmt.Println("End")

}
