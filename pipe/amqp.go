package pipe

import "github.com/streadway/amqp"

type AMQPExchange struct {
	name, kind                            string
	durable, autoDelete, internal, noWait bool
	args                                  amqp.Table
}

type AMQPQueue struct {
	name                                   string
	durable, autoDelete, exclusive, noWait bool
	args                                   amqp.Table
}

type AMQPQueueBinding struct {
	name, key, exchange string
	noWait bool
	args amqp.Table
}

type AMQPPublishing struct {
	exchange, key string
	mandatory, immediate bool
	//msg *amqp.Publishing
}

type AMQPConsumer struct {
	queue, consumer string
	autoAck, exclusive, noLocal, noWait bool
	args amqp.Table
}

func exchangeDeclare(ch *amqp.Channel, exchange *AMQPExchange) error {
	return ch.ExchangeDeclare(
		exchange.name,
		exchange.kind,
		exchange.durable,
		exchange.autoDelete,
		exchange.internal,
		exchange.noWait,
		exchange.args,
	)
}

func queueDeclare(ch *amqp.Channel, queue *AMQPQueue) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queue.name,
		queue.durable,
		queue.autoDelete,
		queue.exclusive,
		queue.noWait,
		queue.args,
	)
}

func queueBind(ch *amqp.Channel, queueBinding *AMQPQueueBinding) error {
	return ch.QueueBind(
		queueBinding.name,
		queueBinding.key,
		queueBinding.exchange,
		queueBinding.noWait,
		queueBinding.args,
	)
}

func populateAMQP(amqpUrl string, exchange *AMQPExchange, queue *AMQPQueue, queueBinding *AMQPQueueBinding, publishing *AMQPPublishing, consumer *AMQPConsumer) (*Pipe, error) {
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = exchangeDeclare(ch, exchange)
	if err != nil {
		return nil, err
	}

	if queue != nil {
		_, err = queueDeclare(ch, queue)
		if err != nil {
			return nil, err
		}

		if queueBinding != nil {
			err = queueBind(ch, queueBinding)
			if err != nil {
				return nil, err
			}
		}
	}

	return &Pipe{
		Connection: conn,
		Channel: ch,
		Publishing: publishing,
		Consumer: consumer,
	}, nil
}
