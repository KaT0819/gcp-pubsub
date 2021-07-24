package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

var PROJECT_ID = os.Getenv("PROJECT_ID")
var PUBSUB_TOPIC = os.Getenv("PUBSUB_TOPIC")
var PORT = os.Getenv("PORT")

func main() {
	if PROJECT_ID == "" || PUBSUB_TOPIC == "" {
		log.Fatal("環境変数を設定してください。PROJECT_ID, PUBSUB_TOPIC")
	}

	http.HandleFunc("/", indexHandler)
	if PORT == "" {
		PORT = "80"
		log.Printf("デフォルトポートの設定 %s", PORT)
	}

	log.Printf("Listening on port %s", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic(PUBSUB_TOPIC)
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("{message: \"Hello Worrld\"}"),
	})

	msgID, err := res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("messageID: %v", msgID)

	_, err = fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// fmt.Fprintln(w, "")
	// fmt.Fprintln(w, "---------------------------")

	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// }
	// for _, addrs := range addrs {
	// 	networkIp, ok := addrs.(*net.IPNet)

	// 	if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
	// 		fmt.Fprintln(w, networkIp.IP.String())
	// 	}
	// }
}
