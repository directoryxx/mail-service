package worker

import (
	"context"
	"fmt"
	"log"
	"mail/infrastructure"
	"mail/internal/controller"
	"mail/internal/repository"
	"mail/internal/usecase"
	"sync"
	"time"

	// "mail/internal/controller"
	// "mail/internal/repository"
	// "mail/internal/usecase"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func RunWorker() {
	// var ctx context.Context = context.Background()

	log.Println("[INFO] App Mode : Worker")

	// app.Validator = &CustomValidator{validator: validator.New()}

	// log.Println("[INFO] Starting mail Service on port", os.Getenv("APPLICATION_PORT"))

	log.Println("[INFO] Loading Database")
	dbMongo, err := infrastructure.ConnectMongo()

	if err != nil {
		log.Fatalf("Could not initialize Mongo connection using client %s", err)
	}

	// Ping the primary
	if err := dbMongo.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	log.Println("[INFO] MongoDB Connected and Pinged")

	defer dbMongo.Disconnect(context.TODO())

	log.Println("[INFO] Loading Kafka Consumer")
	kafkaConn, err := infrastructure.ConnectKafka()

	if err != nil {
		log.Fatalf("Could not initialize connection to Kafka %s", err)
	}

	defer kafkaConn.Close()

	log.Println("[INFO] Loading SMTP Server")
	smtp, err := infrastructure.ConnectSMTP()
	if err != nil {
		log.Fatalf("Could not initialize connection to SMTP %s", err)
	}

	defer smtp.Close()

	log.Println("[INFO] Loading Repository")
	mailRepo := repository.NewMailRepository(smtp, dbMongo)

	log.Println("[INFO] Loading Usecase")
	mailUsecase := usecase.NewMailUseCase(mailRepo)

	log.Println("[INFO] Loading Controller")
	mailController := controller.NewMailController(mailUsecase)

	var wg sync.WaitGroup

	wg.Add(1)

	go consumerKafka("mail", kafkaConn, mailController, &wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
}

func consumerKafka(topic string, kafkaConsumer *kafka.Consumer, mailController controller.MailController, wg *sync.WaitGroup) {
	kafkaConsumer.SubscribeTopics([]string{topic}, nil)

	defer wg.Done()

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := kafkaConsumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			mailController.SendMail(context.TODO(), string(msg.Value))
		}
	}

	kafkaConsumer.Close()
}
