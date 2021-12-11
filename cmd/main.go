package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"secondMS/repository"
	"secondMS/repository/mongorepo"
	"secondMS/reqhandlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	topic := "add-person"
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	kafkaPort := os.Getenv("KAFKA_PORT")
	//partition := 0
	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoPort := os.Getenv("MONGO_PORT")
	abDb := getDbRepository(mongoAddress, mongoPort)
	repository.InitModule(abDb)
	runKafka(kafkaAddress, kafkaPort, topic, 0)
}

func runKafka(address string, port string, topic string, partition int) {
	address = fmt.Sprintf("%s:%s", address, port)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
		//Partition: partition,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		if string(m.Key) == "Person" {
			reqhandlers.AddPerson(m.Value)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func getDbRepository(address string, port string) *repository.AddressBookRepo {
	address = fmt.Sprintf("mongodb://%s:%s", address, port) // localhost 27017
	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()

	err = client.Connect(ctx)
	if err != nil {
		panic(err.Error())
	}
	db := client.Database("testms")
	abDb := mongorepo.GetAddressBookMongoRepo(db.Collection("AddressBook"))
	return &abDb
}
