package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"go.dataddo.com/pgq/x/schema"

	_ "github.com/jackc/pgx/v4/stdlib"

	"go.dataddo.com/pgq"
)

var queueName = "order_queue"

func CreateQueue() {

	create := schema.GenerateCreateTableQuery(queueName)
	fmt.Println(create)
}

func Publish(orderNumber string) {
	uri := os.Getenv("DATABASE_URI")

	// create a new postgres connection
	db, err := sql.Open("pgx", uri)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close database connection: %v\n", closeErr)
		}
	}()

	// create the publisher which may be reused for multiple messages
	// you may pass the optional PublisherOptions when creating it
	publisher := pgq.NewPublisher(db)

	// publish the message to the queue
	// provide the payload which is the JSON object
	// and optional metadata which is the map[string]string
	message := fmt.Sprintf(`{"order_id":%s}`, orderNumber)
	msg := &pgq.MessageOutgoing{
		Payload: json.RawMessage(message),
	}
	msgID, err := publisher.Publish(context.Background(), queueName, msg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Message published with ID:", msgID)
}

func Consumer() {
	uri := os.Getenv("DATABASE_URI")

	// create a new postgres connection and publisher
	db, err := sql.Open("pgx", uri)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close database connection: %v\n", closeErr)
		}
	}()

	// create the consumer which gets attached to handling function we defined above
	h := &myhandler{}
	consumer, err := pgq.NewConsumer(db, queueName, h)
	if err != nil {
		panic(err.Error())
	}

	err = consumer.Run(context.Background())
	if err != nil {
		panic(err.Error())
	}
}

// we must specify the message handler, which implements simple interface
type myhandler struct{}

func (h *myhandler) HandleMessage(_ context.Context, msg *pgq.MessageIncoming) (processed bool, err error) {
	fmt.Println("Message payload:", string(msg.Payload))
	return true, nil
}
