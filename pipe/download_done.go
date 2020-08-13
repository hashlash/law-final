package pipe

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func downloadDoneRoutingKey(clientId string) string {
	return fmt.Sprintf("download.compress.%s", clientId)
}

func SetupDownloadDonePipeIn(clientId string) (*Pipe, error) {
	amqpUrl := os.Getenv("AMQP_URL")
	exchange := &AMQPExchange{
		name:    os.Getenv("EXCHANGE_DIRECT"),
		kind:    amqp.ExchangeDirect,
		durable: true,
	}
	publishing := &AMQPPublishing{
		exchange: exchange.name,
		key: downloadDoneRoutingKey(clientId),
		mandatory: true,
	}
	return populateAMQP(amqpUrl, exchange, nil, nil, publishing, nil)
}

func SetupDownloadDonePipeOut(clientId string) (*Pipe, error) {
	amqpUrl := os.Getenv("AMQP_URL")
	exchange := &AMQPExchange{
		name:    os.Getenv("EXCHANGE_DIRECT"),
		kind:    amqp.ExchangeDirect,
		durable: true,
	}
	queue := &AMQPQueue{
		name:    fmt.Sprintf("server_download_%s_done_q", clientId),
		autoDelete: true,
		exclusive: true,
	}
	queueBinding := &AMQPQueueBinding{
		name: queue.name,
		key: downloadDoneRoutingKey(clientId),
		exchange: exchange.name,
	}
	consumer := &AMQPConsumer{
		queue: queue.name,
		autoAck: true,
		exclusive: true,
	}
	return populateAMQP(amqpUrl, exchange, queue, queueBinding, nil, consumer)
}
