package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strings"
)

var (
	addr    = flag.String("addr", ":8080", "The address to bind to")
	brokers = flag.String("brokers", "0.0.0.0:19092", "The Kafka brokers to connect to, as a comma separated list")
	version = flag.String("version", "2.0.0", "Kafka cluster version")
)

func init() {
	flag.Parse()
}

func ToJSON(w http.ResponseWriter, v interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	respBody := &map[string]string{
		"data": "Hello World",
	}

	// encode response body at the end
	if err := ToJSON(w, respBody, http.StatusCreated); err != nil {
		log.Println(err)
	}
}

func main() {
	brokerList := strings.Split(*brokers, ",")
	fmt.Println(brokerList)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	producer := newSyncProducer(brokerList)

	r.Get("/", handleHome)

	r.Post("/produce", func(w http.ResponseWriter, r *http.Request) {
		payload := new(ProduceRequestDTO)
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// do shit here

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "shit",
			Value: sarama.StringEncoder(payload.Message),
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Printf("produce message to partition: %d with offset: %d\n", partition, offset)

		// response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "message has been produced",
		}); err != nil {
			fmt.Println("encode error")
			w.Write([]byte("error"))
			return
		}
	})

	if err := http.ListenAndServe(func() string {
		log.Printf("http server listen to :3000\n")
		return ":3001"
	}(), r); err != nil {
		log.Fatal(err)
	}
}

func newSyncProducer(brokerList []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal // 1
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("new producer failed: %s", err.Error())
	}
	return producer
}
