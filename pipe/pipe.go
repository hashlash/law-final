package pipe

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Pipe struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
	Publishing *AMQPPublishing
	Consumer *AMQPConsumer
}

func (pipe *Pipe) Close() error {
	err := pipe.Channel.Close()
	if err != nil {
		return err
	}
	return pipe.Connection.Close()
}

func (pipe *Pipe) Send(data map[string]string) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return pipe.Channel.Publish(
		pipe.Publishing.exchange,
		pipe.Publishing.key,
		pipe.Publishing.mandatory,
		pipe.Publishing.immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body: msg,
		},
	)
}

func (pipe *Pipe) Receive() (<-chan map[string]string, error) {
	chMap := make(chan map[string]string)
	chDelivery, err := pipe.Channel.Consume(
		pipe.Consumer.queue,
		pipe.Consumer.consumer,
		pipe.Consumer.autoAck,
		pipe.Consumer.exclusive,
		pipe.Consumer.noLocal,
		pipe.Consumer.noWait,
		pipe.Consumer.args,
	)
	if err != nil {
		return nil, err
	}

	go func() {
		for delivery := range chDelivery {
			data := map[string]string{}
			err := json.Unmarshal(delivery.Body, &data)
			if err != nil {
				log.Println("Error converting delivery body to map: ", err)
				continue
			}
			chMap <- data
		}
	}()

	return (<-chan map[string]string)(chMap), nil
}
