package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()

	// Configura o cliente
	projectid := flag.String("projectid", "", "ID do projeto do Google Cloud")
	subid := flag.String("subid", "", "Subscription Id")
	topicName := flag.String("topic", "golang-conf-01", "Topic Name")
	flag.Parse()

	if *projectid == "" || *subid == "" {
		log.Fatalf("Os argumentos projectID e subscriptionName devem ser fornecidos")
	}

	client, err := pubsub.NewClient(ctx, *projectid)
	if err != nil {
		log.Fatalf("Falha ao criar o cliente PubSub: %v", err)
	}

	// Verifica se a assinatura existe, se não, cria uma nova
	sub := client.Subscription(*subid)
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatalf("Falha ao verificar a existência da assinatura: %v", err)
	}
	if !exists {
		_, err := client.CreateSubscription(ctx, *subid, pubsub.SubscriptionConfig{
			Topic:       client.Topic(*topicName),
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			log.Fatalf("Falha ao criar a assinatura: %v", err)
		}
	}

	pullMessage(ctx, sub)
}

func pullMessage(ctx context.Context, sub *pubsub.Subscription) {
	var mu sync.Mutex
	cctx, _ := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("Recebida mensagem: %q\n", string(msg.Data))
		msg.Ack() // reconhece a mensagem, caso não reconheça, a mensagem voltará para a fila
	})
	if err != nil {
		log.Fatalf("Falha ao receber mensagens: %v", err)
	}
}
