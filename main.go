package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	eventsource "gopkg.in/antage/eventsource.v1"
	"knative.dev/eventing-contrib/pkg/kncloudevents"
	"net/http"
)

/*
Example Output:
☁  cloudevents.Event:
Validation: valid
Context Attributes,
  SpecVersion: 0.2
  Type: dev.knative.eventing.samples.heartbeat
  Source: https://knative.dev/eventing-contrib/cmd/heartbeats/#local/demo
  ID: 3d2b5a1f-10ca-437b-a374-9c49e43c02fb
  Time: 2019-03-14T21:21:29.366002Z
  ContentType: application/json
  Extensions:
    the: 42
    beats: true
    heart: yes
Transport Context,
  URI: /
  Host: localhost:8080
  Method: POST
Data,
  {
    "id":162,
    "label":""
  }
*/
var es eventsource.EventSource

func display(event cloudevents.Event) {
	es.SendEventMessage(event.String(), "tx", "1")
	fmt.Printf("☁️  cloudevents.Event\n%s", event.String())

	//	es.SendEventMessage()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	es = eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/events", es)
	http.HandleFunc("/test", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}

	log.Fatal(c.StartReceiver(context.Background(), display))
	// http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
