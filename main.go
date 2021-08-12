package main

import (
	"context"
	"time"

	"github.com/moooll/cat-service-mongo/internal/handler"
	"github.com/moooll/cat-service-mongo/internal/repository"
	rediscache "github.com/moooll/cat-service-mongo/internal/repository/rediscache"
	service "github.com/moooll/cat-service-mongo/internal/service"
	"github.com/moooll/cat-service-mongo/internal/streams"
	"github.com/segmentio/kafka-go"

	"log"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(repository.DatabaseURI))
	if err != nil {
		log.Print("error connecting to the db\n", err.Error())
	}

	defer func() {
		err = mongoClient.Disconnect(ctx)
		if err != nil {
			log.Print("error disconnecting from the db\n", err.Error())
		}
	}()

	collection := mongoClient.Database("catalog").Collection("cats2")
	dbs, err := mongoClient.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		log.Print("error listing dbs ", err.Error())
	}

	collections, _ := mongoClient.Database("catalog").ListCollectionNames(context.Background(), bson.M{})
	log.Print(collections)
	log.Print(dbs)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
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
	kafkaConn, err := kafkaConnect()
	if err != nil {
		log.Println("error connecting to Kafka ", err.Error())
	}

	go func() {
		for {
			data, err := ss.Read(ctx, "$")
			if err != nil {
				log.Println("error listening on delete: ", err.Error())
			}
		
			_, err = kafkaConn.WriteMessages(kafka.Message{
				Key:   []byte("delete-cats:"),
				Value: data.([]byte),
				Time:  time.Now(),
			})
			if err != nil {
				log.Println("error writing to Kafka: ", err.Error())
			}
		}
	}()

	batch := kafkaConn.ReadBatch(10e3, 10e6)
	
	defer func() {
		if err = batch.Close(); err != nil {
			log.Println("error closing Kafka batch ", err.Error())
		}

		if err = kafkaConn.Close(); err != nil {
			log.Println("error closing Kafka connection ", err.Error())
		}
	}()
		
	b := make([]byte, 10e3)

	go func() {
		for {
			err = service.ReadFromKafka(ctx, kafkaConn, batch, b)
			if err != nil {
				log.Println("error reading messages from Kafka ", err.Error())
			}
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

func kafkaConnect() (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), ":9092", "tcp", "delete-cats", 0)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
