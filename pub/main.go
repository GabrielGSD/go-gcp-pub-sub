package main

import (
	"context"
	"flag"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	projectid := flag.String("projectid", "go-conf-417903", "Project ID")
	topicName := flag.String("topic", "golang-conf-01", "Topic Name")
	msg := flag.String("msg", "hi", "Message")
	flag.Parse()

	client, err := pubsub.NewClient(ctx, *projectid)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	topic := client.Topic(*topicName)
	exist, err := topic.Exists(ctx)
	if err != nil {
		panic(err)
	}

	if !exist {
		topic, err = client.CreateTopic(ctx, *topicName)
		if err != nil {
			panic(err)
		}
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(*msg),
	})

	_, err = result.Get(ctx)
	if err != nil {
		panic(err)
	}

}
