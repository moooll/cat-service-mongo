package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/moooll/cat-service-mongo/internal"
	"github.com/moooll/cat-service-mongo/internal/handler"
	"github.com/moooll/cat-service-mongo/internal/repository"
	rediscache "github.com/moooll/cat-service-mongo/internal/repository/rediscache"
	service "github.com/moooll/cat-service-mongo/internal/service"
	"github.com/moooll/cat-service-mongo/internal/streams"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/ffjson/ffjson"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(internal.MongoURI))
	if err != nil {
		log.Errorln("error connecting to the db ", err.Error())
	}

	defer func() {
		err = mongoClient.Disconnect(mongoCtx)
		if err != nil {
			log.Errorln("error disconnecting from the db ", err.Error())
		}
	}()

	collection := mongoClient.Database("catalog").Collection("cats2")
	dbs, err := mongoClient.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		log.Errorln("error listing dbs ", err.Error())
	}

	collections, err := mongoClient.Database("catalog").ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Errorln("error listing collections ", err.Error())
	}

	log.Infoln(collections)
	log.Infoln(dbs)

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
	kafkaWriter := kafkaW()
	kafkaReader := kafkaR()

	defer func() {
		if err = kafkaReader.Close(); err != nil {
			log.Errorln("error closing Kafka Reader: ", err.Error())
		}

		if err = kafkaWriter.Close(); err != nil {
			log.Errorln("error closing Kafka Writer: ", err.Error())
		}
	}()

	go writeFromRedisToKafka(ss, kafkaWriter)
	go readFromKafka(kafkaReader)
	err = echoStart(serv)
	if err != nil {
		log.Errorln("could not start server ", err.Error())
	}
}

func kafkaW() *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(internal.KafkaURI),
		Topic: "delete-cats",
	}
}

func kafkaR() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{internal.KafkaURI},
		Topic:     "delete-cats",
		Partition: 0,
		MaxBytes:  10e3,
		MinBytes:  100,
	})
}

func echoStart(serv *handler.Service) error {
	e := echo.New()
	e.POST("/cats", serv.AddCat)
	e.GET("/cats", serv.GetAllCats)
	e.GET("/cats/:id", serv.GetCat)
	e.PUT("/cats", serv.UpdateCat)
	e.DELETE("/cats/:id", serv.DeleteCat)
	if err := e.Start(":8081"); err != nil {
		return err
	}
	return nil
}

func writeFromRedisToKafka(ss *streams.StreamService, w *kafka.Writer) {
	i := 0
	for {
		data, er := ss.Read(context.Background(), "$")
		if er != nil {
			log.Errorln("error reading from Redis stream: ", er.Error())
		}

		dataB, e := ffjson.Marshal(&data)
		if e != nil {
			log.Errorln("error marshaling data from redis stream: ", e.Error())
		}
		err := w.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprint(i)),
			Value: dataB,
		})
		if err != nil {
			log.Errorln("error writing to Kafka: ", err.Error())
		}
		i++
	}
}

func readFromKafka(r *kafka.Reader) {
	key := []byte("")
	for {
		mes, errr := r.ReadMessage(context.Background())
		if errr != nil {
			log.Errorln("error reading from Kafka: ", errr.Error())
		}

		if !bytes.Equal(mes.Key, key) {
			log.Infoln("Kafka message is: ", string(mes.Value))
		}

		key = mes.Key
	}
}
