package pipe

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func DownloadProgressRoutingKey(clientId string) string {
	return fmt.Sprintf("client.%s.download", clientId)
}

func SetupDownloadProgressPipeIn(clientId string) (*Pipe, error) {
	amqpUrl := os.Getenv("AMQP_URL")
	exchange := &AMQPExchange{
		name:    os.Getenv("EXCHANGE_TOPIC"),
		kind:    amqp.ExchangeTopic,
		durable: true,
	}
	publishing := &AMQPPublishing{
		exchange: exchange.name,
		key: DownloadProgressRoutingKey(clientId),
		mandatory: true,
	}
	return populateAMQP(amqpUrl, exchange, nil, nil, publishing, nil)
}
