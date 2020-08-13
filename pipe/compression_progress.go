package pipe

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func CompressProgressRoutingKey(clientId string) string {
	return fmt.Sprintf("client.%s.compress", clientId)
}

func SetupCompressionProgressQueue(clientId string) (*Pipe, error) {
	amqpUrl := os.Getenv("AMQP_URL")
	exchange := &AMQPExchange{
		name:    os.Getenv("EXCHANGE_TOPIC"),
		kind:    amqp.ExchangeTopic,
		durable: true,
	}
	publishing := &AMQPPublishing{
		exchange: exchange.name,
		key: CompressProgressRoutingKey(clientId),
		mandatory: true,
	}
	return populateAMQP(amqpUrl, exchange, nil, nil, publishing, nil)
}
