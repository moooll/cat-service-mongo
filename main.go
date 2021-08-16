package main

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/moooll/cat-service-mongo/internal"
	"github.com/moooll/cat-service-mongo/internal/handler"
	"github.com/moooll/cat-service-mongo/internal/repository"
	rediscache "github.com/moooll/cat-service-mongo/internal/repository/rediscache"
	service "github.com/moooll/cat-service-mongo/internal/service"
	"github.com/moooll/cat-service-mongo/internal/streams"

	"log"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(internal.MongoURI))
	if err != nil {
		log.Print("error connecting to the db\n", err.Error())
	}

	defer func() {
		err = mongoClient.Disconnect(mongoCtx)
		if err != nil {
			log.Print("error disconnecting from the db\n", err.Error())
		}
	}()

	collection := mongoClient.Database("catalog").Collection("cats2")
	dbs, err := mongoClient.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		log.Print("error listing dbs ", err.Error())
	}

	collections, err := mongoClient.Database("catalog").ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Print("error listing collections ", err.Error())
	}

	log.Print(collections)
	log.Print(dbs)

	rdb := redis.NewClient(&redis.Options{
		Addr:     internal.RedisURI,
		Password: "",
		DB:       0,
	})

	redisC := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	ss := streams.NewStreamService(rdb)
	serv := handler.NewService(
		service.NewStorage(repository.NewCatalog(collection), rediscache.NewRedisCache(redisC, rdb)), ss)

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{internal.KafkaURI},
		Topic:     "delete-cats",
		Partition: 0,
		MaxBytes:  10e6,
		MinBytes:  10e3,
	})

	kafkaWriter := &kafka.Writer{
		Addr:  kafka.TCP(internal.KafkaURI),
		Topic: "delete-cats",
	}
	if err != nil {
		log.Println("error connecting to Kafka ", err.Error())
	}

	defer func() {
		if err = kafkaReader.Close(); err != nil {
			log.Println("error closing Kafka Reader: ", err.Error())
		}

		if err = kafkaWriter.Close(); err != nil {
			log.Println("error closing Kafka Writer: ", err.Error())
		}
	}()

	go func() {
		for {
			data, err := ss.Read(context.Background(), "$")
			if err != nil {
				log.Println("error reading from Redis stream: ", err.Error())
			}

			dataB, err := ffjson.Marshal(&data)
			if err != nil {
				log.Println("error marshaling data from redis stream: ", err.Error())
			}
			err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte("delete-cats:"),
				Value: dataB,
			})
			if err != nil {
				log.Println("error writing to Kafka: ", err.Error())
			}
		}
	}()

	go func() {
		for {
			mes, err := kafkaReader.ReadMessage(context.Background())
			if err != nil {
				log.Println("error reading from Kafka: ", err.Error())
			}

			log.Println("Kafka message: ", string(mes.Value))
		}
	}()

	e := echo.New()
	e.POST("/cats", serv.AddCat)
	e.GET("/cats", serv.GetAllCats)
	e.GET("/cats/:id", serv.GetCat)
	e.PUT("/cats", serv.UpdateCat)
	e.DELETE("/cats/:id", serv.DeleteCat)
	e.GET("/cats/get-rand-cat", handler.GetRandCat)
	if err := e.Start(":8081"); err != nil {
		log.Print("could not start server\n", err.Error())
	}
}
