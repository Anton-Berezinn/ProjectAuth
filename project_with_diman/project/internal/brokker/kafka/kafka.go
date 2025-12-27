package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	baseUrl = "http://localhost:8083/get"
)

func NewClient(data []byte) {

	params := url.Values{}
	params.Add("data", string(data))

	fullURL := baseUrl + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела:", err)
		return
	}

	fmt.Println("Ответ сервера:", string(body))
}

func Consumer() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
		GroupID: "my-groupID",
	})
	defer reader.Close()

	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Println("Error to ReadMessage", err.Error())
		return err
	}

	//Todo: Можно сделать клиента который будет отправлять запрос на другой сервер
	fmt.Println(string(msg.Value))
	NewClient(msg.Value)
	return nil

}

func Producer(ctx context.Context, email string, userID int) error {

	// At least once гарантируем что сообщение будет отправлено в kafky
	// будет RequiredAcks: -1  гарантия доставки
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "my-topic",
		Balancer: &kafka.RoundRobin{},
	})

	defer writer.Close()

	var value strings.Builder

	value.WriteString(fmt.Sprintf(" User connect to serve with email %s and UserId %d", email, userID))

	err := writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(value.String()),
	})
	//err := writer.WriteMessages(context.Background(), kafka.Message{
	//		Key:   []byte("user-123"), //устанавливаем ключ
	//		Value: []byte("Hello, Kafka!"),
	//	}) В этом примере все сообщения с ключом "user-123" попадут в одну и ту же партицию, обеспечивая порядок обработки. Это сделано для маршрутизации.
	if err != nil {
		log.Println("Error writing to kafka", err)
		return err
	}
	fmt.Println("Сообщение отправлено с гарантией At least once.")
	Consumer()
	return nil
}
