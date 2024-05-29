package message_broker

import (
	"fmt"
	"go_test_rabbitmq/config"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	Conn 	*amqp091.Connection
	Ch 		*amqp091.Channel
}


type Options struct {
	ExchangeDurable , ExchangeAutoDelete , ExchangeInternal , ExchangeNoWait bool
	Args amqp091.Table
	QueueDurable , QueueAutoDelete , QueueExclusive , QueueNoWait bool
	BindNoWait bool
}

var DefaultOptions = &Options{
	ExchangeDurable: true,
	QueueDurable: true,
	ExchangeAutoDelete: false,
	QueueAutoDelete: false,
	ExchangeInternal: false,//for alternate exchange
	QueueExclusive: false,// one queue for one consume,
	ExchangeNoWait: false,
	QueueNoWait: false,
	BindNoWait: false,
	Args: nil,
}

func (rc *RabbitMQClient) NewExchangeQueue(exchangeName , queueName , routingKey , typeExcange string , options *Options) error {
	if typeExcange == amqp091.ExchangeDirect {
		if exchangeName != routingKey {
			return fmt.Errorf("routing key '%s' exchange name  '%s' must be same because exchange direct",routingKey,exchangeName)
		}
	}
	if options == nil {
		options = DefaultOptions
	}
	fmt.Println("Processing ...")
	err := rc.Ch.ExchangeDeclare(exchangeName,typeExcange,options.ExchangeDurable,options.ExchangeAutoDelete,options.ExchangeInternal,options.ExchangeNoWait,options.Args)
	if err != nil {
		return nil
	}
	fmt.Println("Successfully crete exchange and processing create queue")
	queue , err := rc.Ch.QueueDeclare(queueName,options.QueueDurable,options.QueueAutoDelete,options.QueueExclusive,options.QueueNoWait,options.Args)
	if err != nil {
		return nil
	}
	fmt.Println("Successfully crete queue and processing create bind")
	err = rc.Ch.QueueBind(queue.Name,routingKey,exchangeName,options.BindNoWait,options.Args)
	if err != nil {
		return nil
	}
	return nil
}

func (rc *RabbitMQClient) Close() error {
	return rc.Ch.Close()
}

func NewRabbitMQClient(con *amqp091.Connection) (*RabbitMQClient,error) {
	ch  , err := con.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{
		Conn: con,
		Ch: ch,
	}, nil
}



func ConnectionRabbitMQ(url string)  (*amqp091.Connection, error) {
	// protocol://user:password@host:port/virtualhost
	if url == "" {
		cfg := config.NewAppConfig()
		url = fmt.Sprintf("%s://%s:%s@%s/%s",cfg.GetString("rabbitmq.protocol"),cfg.GetString("rabbitmq.user"),cfg.GetString("rabbitmq.password"),cfg.GetString("rabbitmq.host"),cfg.GetString("rabbitmq.vhost"))
	}	
	return amqp091.Dial(url)
}