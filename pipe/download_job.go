package pipe

import (
	"github.com/streadway/amqp"
	"os"
)

func SetupDownloadJobPipe() (*Pipe, error) {
	amqpUrl := os.Getenv("AMQP_URL")
	exchange := &AMQPExchange{
		name:    os.Getenv("EXCHANGE_DIRECT"),
		kind:    amqp.ExchangeDirect,
		durable: true,
	}
	queue := &AMQPQueue{
		name:    "server_download_q",
		durable: true,
	}
	queueBinding := &AMQPQueueBinding{
		name: queue.name,
		key: "download.server",
		exchange: exchange.name,
	}
	publishing := &AMQPPublishing{
		exchange: exchange.name,
		key: queueBinding.key,
		mandatory: true,
	}
	consumer := &AMQPConsumer{
		queue: queue.name,
		autoAck: true,
	}

	return populateAMQP(
		amqpUrl, exchange, queue, queueBinding, publishing, consumer)
}
